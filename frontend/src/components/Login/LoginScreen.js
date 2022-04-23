import { botName } from "../../const";
import TelegramLoginButton from "react-telegram-login";

import styles from "./LoginScreen.module.css";

function LoginScreen() {
  const redirectUrl = window.location.href;
  return (
      <TelegramLoginButton
        className={styles.centered}
        dataAuthUrl={redirectUrl}
        botName={botName}
      />
  );
}

export default LoginScreen;
