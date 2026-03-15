import { useNavigate } from "react-router-dom";
import useAuth from "../hooks/useAuth";
import { useState } from "react";
import {
    Container,
    Box,
    Typography,
    TextField,
    Button,
    Alert,
    CircularProgress,
    Paper,
    FormHelperText,
    Card,
    CardContent,
} from "@mui/material";
import { Upload as UploadIcon } from "lucide-react";

interface UploadFormData {
    modelName: string;
    file: File | null;
}

export default function UploadModel() {
    const navigate = useNavigate();
    const { user, getToken } = useAuth();
    const [formData, setFormData] = useState<UploadFormData>({
        modelName: "",
        file: null,
    });
    const [error, setError] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);
    const [uploading, setUploading] = useState(false);

    const handleModelNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData((prev) => ({
            ...prev,
            modelName: e.target.value,
        }));
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            // Validate file type
            if (!file.name.endsWith(".pkl") && file.type !== "application/octet-stream") {
                setError("Please select a valid .pkl file");
                return;
            }
            setFormData((prev) => ({
                ...prev,
                file: file,
            }));
            setError(null);
        }
    };

    const handleUpload = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError(null);
        setSuccessMessage(null);

        // Validation
        if (!formData.modelName.trim()) {
            setError("Model name is required");
            return;
        }

        if (!formData.file) {
            setError("Please select a file to upload");
            return;
        }

        setUploading(true);

        try {
            const token = await getToken();
            if (!token || !user) {
                setError("You must be logged in to upload a model.");
                return;
            }

            // Create FormData for multipart upload
            const uploadFormData = new FormData();
            uploadFormData.append("modelName", formData.modelName);
            uploadFormData.append("file", formData.file);

            // Upload to backend
            const uploadResponse = await fetch("/models/submitModel", {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${token}`,
                },
                body: uploadFormData,
            });

            if (!uploadResponse.ok) {
                const errorData = await uploadResponse.json().catch(() => ({}));
                throw new Error(errorData.error || `Upload failed with status ${uploadResponse.status}`);
            }

            const responseData = await uploadResponse.json();

            setSuccessMessage(
                responseData.message || "Model uploaded successfully! Check your models page to manage it."
            );
            
            // Reset form
            setFormData({
                modelName: "",
                file: null,
            });

            // Redirect to manage models after a short delay
            setTimeout(() => {
                navigate("/models");
            }, 2000);
        } catch (err: any) {
            setError(err.message || "Failed to upload model");
        } finally {
            setUploading(false);
        }
    };

    return (
        <Container maxWidth="md" sx={{ py: 4 }}>
            <Box sx={{ mb: 4 }}>
                <Typography variant="h4" component="h1" sx={{ fontWeight: 700, mb: 1 }}>
                    Upload Your Model
                </Typography>
                <Typography variant="body1" color="textSecondary">
                    Upload your trained .pkl (pickle) machine learning model to share with the
                    prediction pool.
                </Typography>
            </Box>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error}
                </Alert>
            )}

            {successMessage && (
                <Alert severity="success" sx={{ mb: 3 }}>
                    {successMessage}
                </Alert>
            )}

            <Card>
                <CardContent>
                    <form onSubmit={handleUpload}>
                        <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
                            <TextField
                                label="Model Name"
                                placeholder="e.g., My Baseball Predictor v1"
                                fullWidth
                                value={formData.modelName}
                                onChange={handleModelNameChange}
                                disabled={uploading}
                                helperText="Choose a descriptive name for your model"
                            />

                            <Box>
                                <Typography variant="subtitle2" sx={{ mb: 1, fontWeight: 600 }}>
                                    Model File (.pkl)
                                </Typography>
                                <Paper
                                    variant="outlined"
                                    sx={{
                                        p: 3,
                                        textAlign: "center",
                                        cursor: uploading ? "not-allowed" : "pointer",
                                        backgroundColor: formData.file ? "#f5f5f5" : "transparent",
                                        border: formData.file
                                            ? "2px solid #1976d2"
                                            : "2px dashed #ccc",
                                        transition: "all 0.3s ease",
                                        "&:hover": {
                                            borderColor: uploading ? "#ccc" : "#1976d2",
                                            backgroundColor: uploading ? "transparent" : "#fafafa",
                                        },
                                    }}
                                    component="label"
                                >
                                    <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 1 }}>
                                        <UploadIcon size={40} color="#1976d2" />
                                        <Typography variant="subtitle1">
                                            {formData.file ? formData.file.name : "Click to select or drag and drop"}
                                        </Typography>
                                        <Typography variant="caption" color="textSecondary">
                                            {formData.file
                                                ? `${(formData.file.size / 1024 / 1024).toFixed(2)} MB`
                                                : "Only .pkl files are supported"}
                                        </Typography>
                                    </Box>
                                    <input
                                        type="file"
                                        hidden
                                        accept=".pkl"
                                        onChange={handleFileChange}
                                        disabled={uploading}
                                    />
                                </Paper>
                                <FormHelperText sx={{ mt: 1 }}>
                                    Upload a Python pickle file (.pkl) containing your trained model
                                </FormHelperText>
                            </Box>

                            <Box sx={{ display: "flex", gap: 2, justifyContent: "flex-end" }}>
                                <Button
                                    variant="outlined"
                                    onClick={() => navigate("/models")}
                                    disabled={uploading}
                                >
                                    Cancel
                                </Button>
                                <Button
                                    variant="contained"
                                    type="submit"
                                    disabled={uploading || !formData.modelName.trim() || !formData.file}
                                    sx={{ minWidth: 120 }}
                                >
                                    {uploading ? (
                                        <>
                                            <CircularProgress size={20} sx={{ mr: 1 }} />
                                            Uploading...
                                        </>
                                    ) : (
                                        "Upload Model"
                                    )}
                                </Button>
                            </Box>
                        </Box>
                    </form>
                </CardContent>
            </Card>

            <Box sx={{ mt: 4, p: 2, backgroundColor: "#f5f5f5", borderRadius: 1 }}>
                <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1 }}>
                    📝 Tips for uploading models:
                </Typography>
                <ul style={{ margin: 0, paddingLeft: 20 }}>
                    <li>Ensure your model is saved as a .pkl (pickle) file</li>
                    <li>Use descriptive model names to help you identify them later</li>
                    <li>Keep file sizes reasonable for optimal performance</li>
                    <li>You can upload multiple versions of your models</li>
                </ul>
            </Box>
        </Container>
    );
}