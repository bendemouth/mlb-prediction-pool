import { useNavigate } from "react-router-dom";
import { signIn } from "aws-amplify/auth";
import { useState } from "react";
import { Container, Box, Typography, TextField, Button, Alert } from "@mui/material";

function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();

    const handleLogin = async (e: React.SubmitEvent) => {
        e.preventDefault();
        try {
            await signIn({ username: email, password });
            navigate("/");
        } catch (err) {
            setError("Invalid email or password");
        }
    };

    return (
        <Container maxWidth="sm">
            <Box component="form" onSubmit={handleLogin} sx={{ mt: 8 }}>
                <Typography variant="h4" color="primary" gutterBottom>
                    Sign In
                </Typography>
                {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
                <TextField
                    fullWidth label="Email" type="email" margin="normal"
                    value={email} onChange={e => setEmail(e.target.value)}
                />
                <TextField
                    fullWidth label="Password" type="password" margin="normal"
                    value={password} onChange={e => setPassword(e.target.value)}
                />
                <Button type="submit" variant="contained" fullWidth sx={{ mt: 2 }}>
                    Sign In
                </Button>
            </Box>
        </Container>
    );
}

export default Login;