import styles from "./Item.module.css";

import AddItemControls from "./AddControls";

function Item(props) {
  const updateQuantity = (quantity) => {
    props.updateQuantity(props.name, quantity);
  };

  return (
    <div className={styles.item}>
      <div className={styles["item-name"]}>{props.name}</div>
      {props.quantity > 0 ? (
        <div className={styles["item-quantity"]}>{props.quantity}</div>
      ) : null}
      <AddItemControls
        updateQuantity={updateQuantity}
        quantity={props.quantity}
      />
    </div>
  );
}

export default Item;
