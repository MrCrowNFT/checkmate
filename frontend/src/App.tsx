import "./App.css";
import Auth from "./components/auth";
import Deployents from "./components/deployments";
import { AuthProvider, useAuth } from "./contexts/auth-context";

function AuthenticatedApp() {
  const { currentUser } = useAuth();

  return <div>{currentUser ? <Deployents /> : <Auth />}</div>;
}

function App() {
  return (
    <>
      <AuthProvider>
        <AuthenticatedApp />
      </AuthProvider>
    </>
  );
}

export default App;
