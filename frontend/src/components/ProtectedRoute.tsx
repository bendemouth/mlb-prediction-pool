import useAuth from "../hooks/useAuth";
import CircularProgress from "@mui/material/CircularProgress";
import { Navigate, useLocation } from "react-router";
import { Box } from "@mui/material";
import React from "react";

type ProtectedRouteProps = {
    children: React.ReactNode;
    allowWithoutProfile?: boolean;
}

function ProtectedRoute({ children, allowWithoutProfile = false }: ProtectedRouteProps) {
    const { isAuthenticated, hasProfile, loading } = useAuth();
    const currentLocation = useLocation();

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 8}}>
                <CircularProgress />
            </Box>
        )
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" state={{ from: currentLocation }} replace />;
    }

    if (hasProfile === false && currentLocation.pathname !== "/setup-profile") {
        return <Navigate to="/setup-profile" state={{ from: currentLocation }} replace />;
    }

    return <>{children}</>;
}

export default ProtectedRoute;