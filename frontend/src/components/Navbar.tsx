import * as React from "react";
import { NavLink } from "react-router-dom";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import MenuIcon from "@mui/icons-material/Menu";

export type NavItem = {
    label: string;
    to: string;
};

type AppNavbarProps = {
    title?: React.ReactNode;
    navItems?: NavItem[];
}

const defaultNavItems: NavItem[] = [
    { label: "Home", to: "/" },
    { label: "Profile", to: "/profile" },
    { label: "Leaderboard", to: "/leaderboard" },
    { label: "Predictions", to: "/predictions" },
];

export default function AppNavbar(props: AppNavbarProps) {
    const { title = "(ML)B Predictions", navItems = defaultNavItems } = props;
    const [anchorElNav, setAnchorElNav] = React.useState<null | HTMLElement>(null);

    const handleOpenMenu = (event: React.MouseEvent<HTMLElement>) => 
        setAnchorElNav(event.currentTarget);

    const handleCloseMenu = () => setAnchorElNav(null);

    return (
        <AppBar position="sticky">
            <Container maxWidth="lg">
                <Toolbar disableGutters sx={{ gap: 2}}>
                    <Typography
                        variant="h6"
                        component={NavLink}
                        to="/"
                        sx={{
                            mr: 2,
                            fontWeight: 800,
                            color: "inherit",
                            textDecoration: "none"
                        }}
                    >
                        {title}
                    </Typography>
                    <Box sx={{ flexGrow: 1}} />
                    {/* Mobile menu */}
                    <Box sx={{ display: { xs: "flex", md: "none" }, gap: 2 }}>
                        <IconButton
                            size="large"
                            aria-label="open navigation menu"
                            onClick={handleOpenMenu}
                            color="inherit"
                            >
                            <MenuIcon />
                        </IconButton>
                        <Menu
                            anchorEl={anchorElNav}
                            open={Boolean(anchorElNav)}
                            onClose={handleCloseMenu}
                            anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
                            transformOrigin={{ vertical: "top", horizontal: "right" }}
                            >
                            {navItems.map((item) => (
                                <MenuItem onClick={handleCloseMenu}>
                                    <NavLink
                                        to={item.to}
                                        style={{ textDecoration: "none", color: "inherit", width: "100%" }}
                                    >
                                        {({ isActive }) => (
                                            <Box sx={{ fontWeight: isActive ? 700 : 500 }}>
                                                {item.label}
                                            </Box>
                                        )}
                                    </NavLink>
                                </MenuItem>
                            ))}
                            </Menu>
                    </Box>
                    {/* Desktop menu */}
                    <Box sx={{ display: { xs: "none", md: "flex" }, gap: 1 }}>
                        {navItems.map((item) => (
                            <Button
                                key={item.to}
                                component={NavLink}
                                to={item.to}
                                sx={{
                                    color: "inherit",
                                    "&.active": {
                                        backgroundColor: theme => theme.palette.primary.main,
                                        color: theme => theme.palette.primary.contrastText,
                                    },
                                    "&:hover": {
                                        backgroundColor: theme => theme.palette.primary.main,
                                        color: theme => theme.palette.primary.contrastText,
                                    },
                                }}>
                                {item.label}
                            </Button>
                        ))}
                    </Box>
                </Toolbar>
            </Container>
        </AppBar>
    )
}