import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { signIn, signUp, confirmSignUp } from "aws-amplify/auth";
import {
  Container, Box, Typography, TextField, Button, Alert, Tabs, Tab,
} from "@mui/material";
import React from "react";

type TabValue = "signin" | "signup";

function Login() {
  const [tab, setTab] = useState<TabValue>("signin");

  // Sign-in state
  const [signInEmail, setSignInEmail] = useState("");
  const [signInPassword, setSignInPassword] = useState("");

  // Sign-up state
  const [signUpEmail, setSignUpEmail] = useState("");
  const [signUpPassword, setSignUpPassword] = useState("");
  const [signUpConfirmPassword, setSignUpConfirmPassword] = useState("");
  const [confirmationCode, setConfirmationCode] = useState("");
  const [needsConfirmation, setNeedsConfirmation] = useState(false);

  const [error, setError] = useState("");
  const [info, setInfo] = useState("");

  const navigate = useNavigate();
  const currentLocation = useLocation();
  const from = (currentLocation.state as { from?: { pathname: string } })?.from?.pathname || "/";

  const handleSignIn = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    try {
      await signIn({ username: signInEmail, password: signInPassword });
      navigate(from, { replace: true });
    } catch (err: any) {
      setError(err.message || "Invalid email or password");
    }
  };

  const handleSignUp = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    if (signUpPassword !== signUpConfirmPassword) {
      setError("Passwords do not match");
      return;
    }
    try {
      await signUp({
        username: signUpEmail,
        password: signUpPassword,
        options: { userAttributes: { email: signUpEmail } },
      });
      setNeedsConfirmation(true);
      setInfo("Check your email for a verification code.");
    } catch (err: any) {
      setError(err.message || "Sign up failed");
    }
  };

  const handleConfirm = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    try {
      await confirmSignUp({ username: signUpEmail, confirmationCode });
      setInfo("Email confirmed! Signing you in...");
      await signIn({ username: signUpEmail, password: signUpPassword });
      navigate("/setup-profile", { replace: true });
    } catch (err: any) {
      setError(err.message || "Confirmation failed");
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 8, display: "flex", flexDirection: "column", gap: 2 }}>
        <Typography variant="h4" textAlign="center">
          (ML)B Predictions
        </Typography>

        <Tabs value={tab} onChange={(_, v) => { setTab(v); setError(""); setInfo(""); }} centered>
          <Tab label="Sign In" value="signin" />
          <Tab label="Create Account" value="signup" />
        </Tabs>

        {error && <Alert severity="error">{error}</Alert>}
        {info && <Alert severity="info">{info}</Alert>}

        {/* --- SIGN IN --- */}
        {tab === "signin" && (
          <Box component="form" onSubmit={handleSignIn} sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <TextField
              label="Email" type="email" value={signInEmail}
              onChange={(e) => setSignInEmail(e.target.value)} required fullWidth
            />
            <TextField
              label="Password" type="password" value={signInPassword}
              onChange={(e) => setSignInPassword(e.target.value)} required fullWidth
            />
            <Button type="submit" variant="contained" size="large" fullWidth>
              Sign In
            </Button>
          </Box>
        )}

        {/* --- SIGN UP --- */}
        {tab === "signup" && !needsConfirmation && (
          <Box component="form" onSubmit={handleSignUp} sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <TextField
              label="Email" type="email" value={signUpEmail}
              onChange={(e) => setSignUpEmail(e.target.value)} required fullWidth
            />
            <TextField
              label="Password" type="password" value={signUpPassword}
              onChange={(e) => setSignUpPassword(e.target.value)} required fullWidth
            />
            <TextField
              label="Confirm Password" type="password" value={signUpConfirmPassword}
              onChange={(e) => setSignUpConfirmPassword(e.target.value)} required fullWidth
            />
            <Button type="submit" variant="contained" size="large" fullWidth>
              Create Account
            </Button>
          </Box>
        )}

        {/* --- EMAIL CONFIRMATION --- */}
        {tab === "signup" && needsConfirmation && (
          <Box component="form" onSubmit={handleConfirm} sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <Typography textAlign="center">
              Enter the verification code sent to <strong>{signUpEmail}</strong>
            </Typography>
            <TextField
              label="Verification Code" value={confirmationCode}
              onChange={(e) => setConfirmationCode(e.target.value)} required fullWidth
            />
            <Button type="submit" variant="contained" size="large" fullWidth>
              Verify & Continue
            </Button>
          </Box>
        )}
      </Box>
    </Container>
  );
}

export default Login;