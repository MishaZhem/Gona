import { Outlet, useNavigate } from "react-router-dom"
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
                <h1>Gona</h1>
                {user && (
                    <button onClick={handleLogout}>
                        Logout
                    </button>
                )}
            </header>
            <div>
                <h1>This is website for Gona</h1>
                <Outlet></Outlet>
            </div>
        </>
    )
}

export default Template
