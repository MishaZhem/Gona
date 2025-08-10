import { Link, Outlet, useNavigate } from "react-router-dom"
import { useAuth } from "../../hooks/useAuth";
import styles from './Template.module.scss'

function Template() {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    return (
        <>
            <header className={styles.header}>
                <Link to="/"><h1>Gona</h1></Link>
                <div className={styles.rightBar}>
                    {user && (
                        <Link to="/profile">
                            <img src={user?.avatar} alt="avatar" className={styles.avatar} />
                        </Link>

                    )}
                    {user && (
                        <button onClick={handleLogout}>
                            Logout
                        </button>
                    )}
                </div>
            </header>
            <Outlet></Outlet>
        </>
    )
}

export default Template
