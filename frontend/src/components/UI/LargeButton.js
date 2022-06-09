import styles from "./LargeButton.module.css";

function LargeButton(props) {
  return (
    <button className={styles.button + " " + props.className} onClick={props.onClick}>
      {props.children}
    </button>
  );
}

export default LargeButton;
