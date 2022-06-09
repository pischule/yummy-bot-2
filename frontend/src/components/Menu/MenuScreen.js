import LargeButton from "../UI/LargeButton";
import ItemList from "./ItemsList";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faDice } from '@fortawesome/free-solid-svg-icons'

import styles from "./Menu.module.css";

function MenuScreen(props) {
  return (
    <div className="menu">
      <h1 className={styles.h1}>{props.title}</h1>
      <ItemList items={props.items} updateQuantity={props.updateQuantity} />
      {props.items.length > 0 && (
        <div className={styles.buttons}>
          <LargeButton className={styles.first} onClick={props.handleRandomClick}>
          <FontAwesomeIcon icon={faDice} />
          </LargeButton>
          <LargeButton onClick={props.handleButtonClick}>Заказать</LargeButton>
        </div>
      )}
    </div>
  );
}

export default MenuScreen;
