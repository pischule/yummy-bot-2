import styles from "./App.module.css";

import { useState, useEffect } from "react";

import { baseUrl } from "./const";

import MenuScreen from "./components/Menu/MenuScreen";
import ConfirmScreen from "./components/Confirm/ConfirmScreen";
import DoneScreen from "./components/DoneScreen/DoneScreen";

function App() {
  const [screen, setScreen] = useState(0);
  const [items, setItems] = useState([]);
  const [title, setTitle] = useState("Меню не доступно");

  useEffect(() => {
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
    const newItems = items.filter((item) => item.quantity > 0);
    if (newItems.length > 0) {
      setScreen(1);
    }
  };

  const switchToDone = () => {
    setScreen(2);
  };

  const orderedItems = items.filter((item) => item.quantity > 0);

  return (
    <div className={styles.app}>
      {screen === 0 && (
        <MenuScreen
          updateQuantity={updateQuantity}
          handleClick={switchToConfirem}
          items={items}
          title={title}
        />
      )}
      {screen === 1 && (
        <ConfirmScreen switchToDone={switchToDone} items={orderedItems} />
      )}
      {screen === 2 && <DoneScreen />}
    </div>
  );
}

export default App;
