import {
    Container,
    Box,
    Typography,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Alert,
    CircularProgress,
    Card,
    CardContent,
    Chip,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    IconButton,
    Tooltip,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import useAuth from "../hooks/useAuth";
import { useEffect, useState } from "react";
import { Plus, Trash2, Download, Eye } from "lucide-react";
import { ModelMetadata } from "../models/model_metadata";

export default function ManageModels() {
    const navigate = useNavigate();
    const { user, getToken } = useAuth();
    const [models, setModels] = useState<ModelMetadata[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedModel, setSelectedModel] = useState<ModelMetadata | null>(null);
    const [deleting, setDeleting] = useState(false);

    useEffect(() => {
        fetchUserModels();
    }, [user?.userId]);

    const fetchUserModels = async () => {
        if (!user?.userId) {
            setError("User not authenticated");
            setLoading(false);
            return;
        }

        try {
            setLoading(true);
            setError(null);
            const token = await getToken();

            if (!token) {
                setError("Authentication failed");
                return;
            }

            const response = await fetch("/models", {
                headers: {
                    "Authorization": `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                throw new Error(`Failed to fetch models: ${response.statusText}`);
            }

            const modelsData = await response.json();
            setModels(modelsData || []);
        } catch (err: any) {
            setError(err.message || "Failed to load models");
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteClick = (model: ModelMetadata) => {
        setSelectedModel(model);
        setDeleteDialogOpen(true);
    };

    const handleDeleteConfirm = async () => {
        if (!selectedModel) return;

        setDeleting(true);
        try {
            const token = await getToken();

            if (!token) {
                setError("Authentication failed");
                return;
            }

            const response = await fetch(`/models/delete/${selectedModel.model_id}`, {
                method: "DELETE",
                headers: {
                    "Authorization": `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.error || `Failed to delete model: ${response.statusText}`);
            }

            // Remove from local state
            setModels((prev) => prev.filter((m) => m.model_id !== selectedModel.model_id));
            setDeleteDialogOpen(false);
            setSelectedModel(null);
        } catch (err: any) {
            setError(err.message || "Failed to delete model");
        } finally {
            setDeleting(false);
        }
    };

    const handleDownload = (model: ModelMetadata) => {
        // TODO: Implement backend endpoint to generate a download URL
        // Expected endpoint: POST /models/{modelId}/download
        // Should return a signed S3 URL for downloading the model file
        console.log("Downloading model:", model.model_id);
    };

    const getStatusColor = (status: string): "default" | "primary" | "secondary" | "error" | "info" | "success" | "warning" => {
        switch (status?.toLowerCase()) {
            case "active":
                return "success";
            case "pending":
                return "info";
            case "failed":
                return "error";
            default:
                return "default";
        }
    };

    const formatDate = (date: string | Date) => {
        const d = new Date(date);
        return d.toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 4, textAlign: "center" }}>
                <CircularProgress />
                <Typography variant="h6" sx={{ mt: 2 }}>
                    Loading your models...
                </Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 4 }}>
                <Box>
                    <Typography variant="h4" component="h1" sx={{ fontWeight: 700, mb: 1 }}>
                        Manage Your Models
                    </Typography>
                    <Typography variant="body1" color="textSecondary">
                        View, manage, and delete your uploaded ML models
                    </Typography>
                </Box>
                <Button
                    variant="contained"
                    startIcon={<Plus size={20} />}
                    onClick={() => navigate("/models/upload")}
                    size="large"
                >
                    Upload New Model
                </Button>
            </Box>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error}
                </Alert>
            )}

            {models.length === 0 ? (
                <Card>
                    <CardContent sx={{ textAlign: "center", py: 6 }}>
                        <Typography variant="h6" color="textSecondary" sx={{ mb: 2 }}>
                            No models uploaded yet
                        </Typography>
                        <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>
                            Start by uploading your first machine learning model to the prediction pool.
                        </Typography>
                        <Button
                            variant="contained"
                            startIcon={<Plus size={20} />}
                            onClick={() => navigate("/models/upload")}
                        >
                            Upload Your First Model
                        </Button>
                    </CardContent>
                </Card>
            ) : (
                <TableContainer component={Paper}>
                    <Table>
                        <TableHead>
                            <TableRow sx={{ backgroundColor: "#f5f5f5" }}>
                                <TableCell sx={{ fontWeight: 600 }}>Model Name</TableCell>
                                <TableCell sx={{ fontWeight: 600 }}>Status</TableCell>
                                <TableCell sx={{ fontWeight: 600 }}>File Name</TableCell>
                                <TableCell sx={{ fontWeight: 600 }}>Created</TableCell>
                                <TableCell sx={{ fontWeight: 600 }}>Last Updated</TableCell>
                                <TableCell align="right" sx={{ fontWeight: 600 }}>
                                    Actions
                                </TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {models.map((model) => (
                                <TableRow key={model.model_id} sx={{ "&:hover": { backgroundColor: "#f9f9f9" } }}>
                                    <TableCell>
                                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                            {model.model_name}
                                        </Typography>
                                    </TableCell>
                                    <TableCell>
                                        <Chip
                                            label={model.status || "Unknown"}
                                            color={getStatusColor(model.status)}
                                            size="small"
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <Typography variant="caption">{model.file_name}</Typography>
                                    </TableCell>
                                    <TableCell>
                                        <Typography variant="caption">
                                            {formatDate(model.created_at)}
                                        </Typography>
                                    </TableCell>
                                    <TableCell>
                                        <Typography variant="caption">
                                            {formatDate(model.updated_at)}
                                        </Typography>
                                    </TableCell>
                                    <TableCell align="right">
                                        <Box sx={{ display: "flex", gap: 0.5, justifyContent: "flex-end" }}>
                                            <Tooltip title="View Details">
                                                <IconButton
                                                    size="small"
                                                    onClick={() => {
                                                        // TODO: Implement view details modal or page
                                                        console.log("View model details:", model.model_id);
                                                    }}
                                                >
                                                    <Eye size={18} />
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Download">
                                                <IconButton
                                                    size="small"
                                                    onClick={() => handleDownload(model)}
                                                >
                                                    <Download size={18} />
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Delete">
                                                <IconButton
                                                    size="small"
                                                    onClick={() => handleDeleteClick(model)}
                                                    sx={{ color: "error.main" }}
                                                >
                                                    <Trash2 size={18} />
                                                </IconButton>
                                            </Tooltip>
                                        </Box>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

            {/* Delete Confirmation Dialog */}
            <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
                <DialogTitle>Delete Model?</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Are you sure you want to delete{" "}
                        <strong>{selectedModel?.model_name}</strong>? This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteDialogOpen(false)} disabled={deleting}>
                        Cancel
                    </Button>
                    <Button
                        onClick={handleDeleteConfirm}
                        color="error"
                        variant="contained"
                        disabled={deleting}
                    >
                        {deleting ? "Deleting..." : "Delete"}
                    </Button>
                </DialogActions>
            </Dialog>
        </Container>
    );
}