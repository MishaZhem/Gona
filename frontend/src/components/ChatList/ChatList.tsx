import { useState } from 'react'
import { chat, close, send } from '../../assets'
import Button from '../Button/Button'
import styles from './ChatList.module.scss'
import { useAuth } from '../../hooks/useAuth'

interface ChatListProps {
  closeChat: () => void
}

function ChatList(props: ChatListProps) {

  return (
    <div>
      <div className={styles.TopChat}>
        <div className={styles.TopName}>
          <div className={styles.TopPhoto}></div>
          <h3>Technical support</h3>
        </div>
        <Button option='invisible' className={styles.closeBtn} onClick={() => props.closeChat()}>
          <img src={close} alt="" />
        </Button>
      </div>
      <div className={styles.messages}>
        <Message
          id="user123"
          text="Hello!"
          name="It's my"
          time={Date.now()}
          icon="/avatars/me.png"
        />
      </div>
    </div>
  )
}

function ChatBottom() {
  return (
    <div className={styles.bottom}>
      <input type="text" className={styles.input} />
      <Button option='accent'>
        <img src={send} alt="" />
      </Button>
    </div>
  )
}

interface MessageProps {
  text: string;
  name: string;
  time: number;
  icon: string;
  id: string;
}

function Message(props: MessageProps) {
  const { user } = useAuth();
  const { text, name, time, icon, id } = props;
  // const MyMessage = id != "bot";
  const MyMessage = user?.id === id;

  const formattedTime = new Date(time).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div className={`${styles.messageWrapper} ${MyMessage ? styles.mine : styles.other}`}>
      <div className={styles.message}>
        {!MyMessage && <img src={icon} className={styles.avatar} alt="avatar" />}
        <div className={styles.content}>
          <div className={styles.meta}>
            <span className={styles.name}>{name}</span>
            <span className={styles.time}>{formattedTime}</span>
          </div>
          <p className={styles.text}>{text}</p>
        </div>
        {MyMessage && <img src={icon} className={styles.avatar} alt="avatar" />}
      </div>
    </div>
  );
}


function Chat() {
  const [opened, setOpened] = useState(false);

  return (
    <>
      {opened ?
        <div className={styles.screenChat}>
          <ChatList closeChat={() => setOpened(false)} />
          <ChatBottom />
        </div> :
        <Button option='accent' round='max' className={styles.openChat} onClick={() => setOpened(true)}>
          <img src={chat} alt="" />
        </Button>}
    </>
  )
}

export default Chat
