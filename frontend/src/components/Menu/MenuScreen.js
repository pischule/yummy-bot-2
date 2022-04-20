import LargeButton from "../UI/LargeButton";
import ItemList from "./ItemsList";

import styles from "./Menu.module.css";

function MenuScreen(props) {
  return (
    <div className="menu">
      <h1 className={styles.h1}>{props.title}</h1>
      <ItemList items={props.items} updateQuantity={props.updateQuantity} />
      {props.items.length > 0 && (
        <LargeButton onClick={props.handleClick}>Заказать</LargeButton>
      )}
    </div>
  );
}

export default MenuScreen;
