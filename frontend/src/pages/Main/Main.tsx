import { Link } from 'react-router-dom'
import { Chat } from '../../components'
import styles from './Main.module.scss'

function Main() {

  return (
    <div className={styles.screen}>
      <Chat></Chat>
      <Link to="/login">Login</Link><br />
      <Link to="/profile">Profile</Link>
    </div>
  )
}

export default Main
