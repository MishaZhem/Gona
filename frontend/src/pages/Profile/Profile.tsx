import { useEffect, useState } from "react";
import { api } from "../../app/axios";

function Profile() {
  const [profile, setProfile] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const getProfile = async () => {
      try {
        const res = await api.get("/profile");
        setProfile(res.data);
      } catch (err: any) {
        console.error("Error:", err.response?.data || err.message);
      } finally {
        setLoading(false);
      }
    };

    getProfile();
  }, []);
  return (
    <div>
      <h1>Hi, {profile?.email || "user"}!</h1>
    </div>
  );
};

export default Profile