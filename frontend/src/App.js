import styles from "./App.module.css";

import { useState, useEffect } from "react";

import MenuScreen from "./components/Menu/MenuScreen";
import ConfirmScreen from "./components/Confirm/ConfirmScreen";
import DoneScreen from "./components/DoneScreen/DoneScreen";
import LoginScreen from "./components/Login/LoginScreen";
import ErrorModal from "./components/UI/ErrorModal";

function App() {
  const [screen, setScreen] = useState("menu");
  const [items, setItems] = useState([]);
  const [title, setTitle] = useState("Меню не доступно");
  const [error, setError] = useState();

  async function fetchData() {
    const result = await fetch(`${process.env.REACT_APP_API_URL}/menu`);
    const json = await result.json();
    let currentId = 0;
    const itemsWithQuantity = json.items.map((item) => {
      return {
        id: currentId++,
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

  const updateQuantity = (id, count) => {
    setItems((prevItems) => {
      const updatedItems = [...prevItems];
      const item = updatedItems.find((item) => item.id === id);
      item.quantity = count;
      return updatedItems;
    });
  };

  const pickRandom = () => {
    setItems((prevItems) => {
      const updatedItems = prevItems.map((item) => {
        item.quantity = 0;
        return item;
      });

      if (prevItems.length < 10) {
        for (let i = 0; i < 2; i++) {
          const randomIndex = Math.floor(Math.random() * prevItems.length);
          updatedItems[randomIndex].quantity = 1;
        }
        return updatedItems;
      }

      const soupsStartIndex = 0;
      const garnishesStartIndex = 4;
      const secondDishesStartIndex = prevItems.length - 6;

      const randomSoupIndex =
        Math.floor(Math.random() * (garnishesStartIndex - soupsStartIndex)) +
        soupsStartIndex;
      const randomGarnishIndex =
        Math.floor(
          Math.random() * (secondDishesStartIndex - garnishesStartIndex)
        ) + garnishesStartIndex;
      const randomSecondDishIndex =
        Math.floor(
          Math.random() * (prevItems.length - secondDishesStartIndex - 1)
        ) + secondDishesStartIndex;

      updatedItems[randomSoupIndex].quantity = 1;
      updatedItems[randomGarnishIndex].quantity = 1;
      updatedItems[randomSecondDishIndex].quantity = 1;

      return updatedItems;
    });
  };

  const errorHandler = () => {
    setError(null);
  };

  const switchToConfirem = () => {
    const selectedItems = items.filter((item) => item.quantity > 0);
    if (selectedItems.length > 0) {
      setScreen("confirm");
    } else {
      setError({
        title: "Заказ пуст",
        message: "выберите хотя бы один пункт",
      });
    }
  };

  const switchToDone = () => {
    setScreen("done");
  };

  const params = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop),
  });
  const id = params.id;
  if (id === null || id === undefined) {
    return (
      <div className={styles.app}>
        <LoginScreen />
      </div>
    );
  }

  const selectedItems = items.filter((item) => item.quantity > 0);

  return (
    <div className={styles.app}>
      {error && (
        <ErrorModal
          title={error.title}
          message={error.message}
          onConfirm={errorHandler}
        />
      )}
      {screen === "menu" && (
        <MenuScreen
          updateQuantity={updateQuantity}
          handleButtonClick={switchToConfirem}
          handleRandomClick={pickRandom}
          items={items}
          title={title}
        />
      )}
      {screen === "confirm" && (
        <ConfirmScreen
          switchToDone={switchToDone}
          items={selectedItems}
          setError={setError}
        />
      )}
      {screen === "done" && <DoneScreen />}
    </div>
  );
}

export default App;
