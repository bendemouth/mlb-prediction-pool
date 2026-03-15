import { Alert, CircularProgress, Container, FormGroup, Typography } from "@mui/material";
import { Upload } from "lucide-react";
import { useState } from "react";

function FormFile(){
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null)

    if (loading){
        return(
            <Container maxWidth="lg"
                sx={{mt: 4, textAlign:"center"}}>
                    <CircularProgress />
                    <Typography variant="h6" sx={{ mt: 2 }}>
                        Loading...
                    </Typography>
                </Container>
        )
    }

    if(error) {
        return (
                <Container maxWidth="lg" sx={{ mt: 4 }}>
                    <Alert severity="error">{error}</Alert>
                </Container>
            );
    }

    return (
        <Container maxWidth="lg" sx={{mt: 4, mb: 4}}>
            <FormGroup>
                
            </FormGroup>
        </Container>
    )
}

export default FormFile