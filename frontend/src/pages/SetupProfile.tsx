import { useState } from "react";
import useAuth from "../hooks/useAuth";
import { useNavigate } from "react-router";
import { Box, Container } from "@mui/system";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Alert from "@mui/material/Alert";
import { Typography } from "@mui/material";
import React from "react";

function SetupProfile() {
    const [username, setUsername] = useState("");
    const [error, setError] = useState("");
    const [submitting, setSubmitting] = useState(false);
    const navigate = useNavigate();
    const { user, getToken, checkAuth } = useAuth();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError("");
        setSubmitting(true);

        try {
            const token = await getToken();
            if (!token || !user) {
                setError("You must be logged in to set up your profile.");
                return;
            }

            const response = await fetch("/users/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify({
                    username,
                    email: user.signInDetails?.loginId ?? ""
                })
            });

            if (!response.ok) {
                const data = await response.json();
                setError(data.error || "Failed to create profile");
            }

            await checkAuth(); // Refresh auth state to update hasProfile
            navigate("/", { replace: true });
        } catch (err: any) {
            setError(err.message || "An error occurred while setting up your profile.");
        } finally {
            setSubmitting(false);
        }
    };

    return (
    <Container maxWidth="sm">
      <Box component="form" onSubmit={handleSubmit} sx={{ mt: 8, display: "flex", flexDirection: "column", gap: 2 }}>
        <Typography variant="h4" textAlign="center">
          Set Up Your Profile
        </Typography>
        <Typography textAlign="center" color="text.secondary">
          Choose a username that will appear on the leaderboard.
        </Typography>
        {error && <Alert severity="error">{error}</Alert>}
        <TextField
          label="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          fullWidth
          inputProps={{ minLength: 3, maxLength: 20 }}
          helperText="3–20 characters, shown publicly on the leaderboard"
        />
        <Button type="submit" variant="contained" size="large" fullWidth disabled={submitting}>
          {submitting ? "Creating Profile..." : "Save & Continue"}
        </Button>
      </Box>
    </Container>
  );
}

export default SetupProfile;