import { useAuth } from "../../hooks/useAuth";
import styles from "./Profile.module.scss";
import { useState, type ChangeEvent } from "react";
import { api } from "../../app/axios";

function Profile() {
  const { user } = useAuth();
  const [selectedTab, setSelectedTab] = useState("settings");
  const [preview, setPreview] = useState<string | null>(user?.avatar_url || null);
  const [_file, setFile] = useState<File | null>(null);
  const [_loading, setLoading] = useState(false);

  const handleAvatarChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0];
    if (!selected) return;

    setFile(selected);
    setPreview(URL.createObjectURL(selected));

    const formData = new FormData();
    formData.append("avatar", selected);

    try {
      setLoading(true);
      await api.post("/upload-avatar", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
    } catch (err) {
      console.error("Avatar upload failed", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div onClick={() => console.log(user)} className={styles.profileWrapper}>
      <aside className={styles.sidebar}>
        <div className={styles.avatarSection}>
          {preview ? (
            <div className={styles.avatarPreview}>
              <img src={preview} alt="Avatar" />
            </div>
          ) : (
            <div className={styles.placeholder}></div>
          )}
          <input type="file" accept="image/*" onChange={handleAvatarChange} />
        </div>

        <nav className={styles.sidebarNav}>
          <button
            onClick={() => setSelectedTab("dashboard")}
            className={selectedTab === "dashboard" ? styles.active : ""}
          >
            Dashboard
          </button>
          <button
            onClick={() => setSelectedTab("messages")}
            className={selectedTab === "messages" ? styles.active : ""}
          >
            Messages
          </button>
          <button
            onClick={() => setSelectedTab("settings")}
            className={selectedTab === "settings" ? styles.active : ""}
          >
            Settings
          </button>
        </nav>
      </aside>

      <main className={styles.settingsPanel}>
        <h2>Settings</h2>
        {selectedTab === "settings" ? (
          <form className={styles.settingsForm}>
            <input type="text" placeholder="Username" defaultValue={user?.name} />
            <input type="email" placeholder="Email" defaultValue={user?.email} />
            <input type="password" placeholder="Current Password" />
            <input type="password" placeholder="New Password" />
            <input type="password" placeholder="Confirm Password" />
            <button type="submit" className={styles.saveBtn}>
              Save Changes
            </button>
          </form>
        ) : (
          <p>This section is under construction</p>
        )}
      </main>
    </div>
  );
}

export default Profile;