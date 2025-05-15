import Auth from "./components/auth";
import Deployments from "./components/deployments";
import { AuthProvider, useAuth } from "./contexts/auth-context";
import { useInitializeTheme } from "./hooks";

function AuthenticatedApp() {
   useInitializeTheme();
  const { currentUser } = useAuth();

  return <div>{currentUser ? <Deployments /> : <Auth />}</div>;
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
