import { useState, useEffect } from "react";

import LargeButton from "../UI/LargeButton";
import ItemList from "./ItemsList";

import styles from "./Menu.module.css";

import { baseUrl } from "../../const";

function MenuScreen(props) {
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

  const handleClick = () => {
    const newItems = items.filter((item) => item.quantity > 0);
    if (newItems.length > 0) {
      setItems(newItems);
      props.switchToConfirm(
        1,
        items.filter((item) => item.quantity > 0)
      );
    }
  };

  return (
    <div className="menu">
      <h1 className={styles.h1}>{title}</h1>
      <ItemList items={items} updateQuantity={updateQuantity} />
      {items.length > 0 && (
        <LargeButton onClick={handleClick}>Заказать</LargeButton>
      )}
    </div>
  );
}

export default MenuScreen;
