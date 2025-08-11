import { Link, Outlet, useNavigate } from "react-router-dom"
import { useAuth } from "../../hooks/useAuth";
import styles from './Template.module.scss'
import { Button } from "../../components";
import { useState } from "react";
import { userLogo } from "../../assets";

function Template() {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [opened, setOpened] = useState<boolean>(false);

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    return (
        <>
            <header className={styles.header}>
                <Link to="/" className={styles.Logo}><h1>Gona</h1></Link>
                <div className={styles.rightBar}>
                    <div className={styles.buttons}>
                        <Link to="/">Catalog</Link>
                        <p>Messages</p>
                        <p>Help</p>
                        <Button>Create new ad</Button>
                        {user && (
                            <div className={styles.profileBlock}>
                                <img onClick={() => setOpened(!opened)} src={user?.avatar_url} alt="avatar" className={styles.avatar} />
                                {opened && (
                                    <div className={styles.menu}>
                                        <Link to="/profile">Profile</Link>
                                        <button onClick={handleLogout}>Logout</button>
                                    </div>
                                )}

                            </div>
                        )}
                        {!user && (
                            <Link to="/login" className={styles.center}>
                                <img src={userLogo} alt="" className={styles.user} />
                            </Link>
                        )}
                    </div>
                </div>
            </header>
            <Outlet></Outlet>
        </>
    )
}

export default Template
