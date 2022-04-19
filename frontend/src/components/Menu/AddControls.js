import styles from "./AddItemControls.module.css";

function AddItemControls(props) {
  const quantity = props.quantity;

  const increment = () => {
    props.updateQuantity(quantity + 1);
  };

  const decrement = () => {
    props.updateQuantity(quantity - 1);
  };

  if (quantity === 0) {
    return (
      <button
        className={`${styles["add-button"]} ${styles.button}`}
        onClick={increment}
      >
        Добавить
      </button>
    );
  } else {
    return (
      <>
        <button
          className={`${styles["pm-button"]} ${styles.button}`}
          onClick={decrement}
        >
          -
        </button>
        <button
          className={`${styles["pm-button"]} ${styles.button}`}
          onClick={increment}
        >
          +
        </button>
      </>
    );
  }
}

export default AddItemControls;
