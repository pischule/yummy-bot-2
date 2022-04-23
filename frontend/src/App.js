import styles from "./App.module.css";

import { useState, useEffect } from "react";

import { baseUrl } from "./const";

import MenuScreen from "./components/Menu/MenuScreen";
import ConfirmScreen from "./components/Confirm/ConfirmScreen";
import DoneScreen from "./components/DoneScreen/DoneScreen";

function App() {
  const [screen, setScreen] = useState("menu");
  const [items, setItems] = useState([]);
  const [title, setTitle] = useState("Меню не доступно");

  async function fetchData() {
    const result = await fetch(`${baseUrl}/menu`);
    const json = await result.json();
    const itemsWithQuantity = json.items.map((item) => {
      return {
        name: item,
        quantity: 0,
      };
    });
    setItems(itemsWithQuantity);
    setTitle(json.title);
  }

  useEffect(() => {
    fetchData();
  }, []);

  const updateQuantity = (name, count) => {
    const newItems = items.map((item) => {
      if (item.name === name) {
        item.quantity = count;
      }
      return item;
    });
    setItems(newItems);
  };

  const switchToConfirem = () => {
    const selectedItems = items.filter((item) => item.quantity > 0);
    if (selectedItems.length > 0) {
      setScreen("confirm");
    }
  };

  const switchToDone = () => {
    setScreen("done");
  };

  const selectedItems = items.filter((item) => item.quantity > 0);

  return (
    <div className={styles.app}>
      {screen === "menu" && (
        <MenuScreen
          updateQuantity={updateQuantity}
          handleButtonClick={switchToConfirem}
          items={items}
          title={title}
        />
      )}
      {screen === "confirm" && (
        <ConfirmScreen switchToDone={switchToDone} items={selectedItems} />
      )}
      {screen === "done" && <DoneScreen />}
    </div>
  );
}

export default App;
