import { send } from '../../assets'
import styles from './ChatList.module.scss'

function ChatList() {

  return (
    <div className={styles.screen}>
      <h1>Hello</h1>
    </div>
  )
}

function ChatBottom() {

  return (
    <div className={styles.bottom}>
      <input type="text" className={styles.input} />
      <button className={styles.send}>
        <img src={send} alt="" />
      </button>
    </div>
  )
}

function Chat() {

  return (
    <div className={styles.screen}>
      <ChatList />
      <ChatBottom />
    </div>
  )
}

export default Chat
