import React, { useState } from "react";
import { useAuth } from "../contexts/auth-context";

const Auth = () => {
  const { signInWithGoogle, signInWithEmail, signUpWithEmail } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLogin, setIsLogin] = useState(true);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleEmailAuth = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      if (isLogin) {
        await signInWithEmail(email, password);
      } else {
        await signUpWithEmail(email, password);
      }
    } catch (error: any) {
      setError(error.message || "Failed to authenticate");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1>Welcome to Checkmate</h1>

      {error && <div>{error}</div>}

      <form onSubmit={handleEmailAuth}>
        <div>
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Processing..." : isLogin ? "Sign In" : "Sign Up"}
        </button>
      </form>

      <p>
        {isLogin ? "Don't have an account? " : "Already have an account? "}
        <button onClick={() => setIsLogin(!isLogin)}>
          {isLogin ? "Sign Up" : "Sign In"}
        </button>
      </p>

      <div>OR</div>

      <button onClick={signInWithGoogle}>Sign in with Google</button>
    </div>
  );
};

export default Auth;
