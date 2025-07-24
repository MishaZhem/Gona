import { useState } from 'react'
import { chat, close, send } from '../../assets'
import Button from '../Button/Button'
import styles from './ChatList.module.scss'

function ChatList() {

  return (
    <div className={styles.screen}>
      <h1>Hello</h1>
    </div>
  )
}

interface ChatBottomProps {
  closeChat: () => void
}

function ChatBottom(props: ChatBottomProps) {

  return (
    <div className={styles.bottom}>
      <input type="text" className={styles.input} />
      <Button option='accent'>
        <img src={send} alt="" />
      </Button>
      <Button option='invisible' className={styles.closeBtn} onClick={() => props.closeChat()}>
        <img src={close} alt="" />
      </Button>
    </div>
  )
}

function Chat() {
  const [opened, setOpened] = useState(false);

  return (
    <>
      {opened ?
        <div className={styles.screenChat}>
          <ChatList />
          <ChatBottom closeChat={() => setOpened(false)} />
        </div> :
        <Button option='accent' round='max' className={styles.openChat} onClick={() => setOpened(true)}>
          <img src={chat} alt="" />
        </Button>}
    </>
  )
}

export default Chat
