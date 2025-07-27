import { useAuth } from "../../hooks/useAuth";

function Profile() {
  const { user } = useAuth();

  return (
    <div>
      <h1>Hi, {user?.email || "user"}!</h1>
    </div>
  );
};

export default Profile