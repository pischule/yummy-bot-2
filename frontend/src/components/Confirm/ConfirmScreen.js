import { useState } from "react";

import styles from "./ConfirmScreen.module.css";

import LargeButton from "../UI/LargeButton";
import Card from "../UI/Card";

import { baseUrl } from "../../const";

function ConfirmScreen(props) {
  const [name, setName] = useState(localStorage.getItem("name") || null);
  const [pressed, setPressed] = useState(false);

  const handleNameChange = (event) => {
    const newName = event.target.value;
    setName(newName);
    localStorage.setItem("name", newName);
  };

  const isNameValid = (name) => {
    const cyrillicPattern = /^[\u0400-\u04FF ]{2,}/;
    return cyrillicPattern.test(name);
  };

  const handleClick = () => {
    if (!isNameValid(name)) {
      alert("Имя может содержать только кириллические символы");
      return;
    }

    if (!pressed && name !== null) {
      setPressed(true);

      const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
      });
      const id = params.id;

      if (id === null || id === undefined) {
        alert("Ошибка авторизации");
        return;
      }

      fetch(`${baseUrl}/order`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          userId: +id,
          name: name,
          items: props.items,
        }),
      }).then((response) => {
        if (response.status === 200) {
          setPressed(false);
          props.switchToDone();
        } else {
          alert("Ошибка отправки заказа");
          setPressed(false);
        }
      });
    }
  };

  return (
    <div>
      <h1 className={styles.h2}>Ваш заказ:</h1>
      <Card>
        <input
          id="nameInput"
          value={name}
          onChange={handleNameChange}
          className={styles.input}
          placeholder="Введите ваше имя"
        ></input>
        <ul>
          {props.items.map((item) => (
            <li key={item.name}>
              {item.name} x{item.quantity}
            </li>
          ))}
        </ul>
      </Card>

      <LargeButton onClick={handleClick}>Отправить</LargeButton>
    </div>
  );
}

export default ConfirmScreen;
