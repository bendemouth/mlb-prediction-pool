import { createTheme } from "@mui/material/styles";

const colors = {
  prussianBlue: "#102542",
  cinnabar: "#e9422f",
  alabasterGrey: "#dadadd",
  khakiBeige: "#b3a394",
  white: "#ffffff",
};

export const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: colors.prussianBlue,
      contrastText: colors.white,
    },
    secondary: {
      main: colors.cinnabar,
    },
    background: {
      default: colors.alabasterGrey,
      paper: colors.white,
    },
    text: {
      primary: colors.prussianBlue,
    },
  },
});