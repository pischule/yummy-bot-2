import TelegramLoginButton from "react-telegram-login";

import styles from "./LoginScreen.module.css";

function LoginScreen() {
  const redirectUrl = window.location.href;
  return (
    <TelegramLoginButton
      className={styles.centered}
      dataAuthUrl={redirectUrl}
      botName={process.env.REACT_APP_BOT_NAME}
    />
  );
}

export default LoginScreen;
