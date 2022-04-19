import styles from "./App.module.css";

import { useState } from "react";

import MenuScreen from "./components/Menu/MenuScreen";
import ConfirmScreen from "./components/Confirm/ConfirmScreen";
import DoneScreen from "./components/DoneScreen/DoneScreen";


function App() {
  const [screen, setScreen] = useState(0);
  const [items, setItems] = useState([]);

  const switchToConfirm = (index, items) => {
    setScreen(1);
    setItems(items);
  }

  const switchToDone = () => {
    setScreen(2);
  }

  return (
    <div className={styles.app}>
      {screen === 0 && <MenuScreen switchToConfirm={switchToConfirm} />}
      {screen === 1 && <ConfirmScreen switchToDone={switchToDone} items={items}/>}
      {screen === 2 && <DoneScreen/>}
    </div>
  );
}

export default App;
