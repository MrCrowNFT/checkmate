import { create } from "zustand";
import { useEffect } from "react";

type ThemeState = {
  darkMode: boolean;
  toggleDarkMode: () => void;
  initializeTheme: () => void;
};

export const useThemeStore = create<ThemeState>((set) => ({
  darkMode:
    typeof window !== "undefined"
      ? localStorage.getItem("darkMode") === "true" ||
        (localStorage.getItem("darkMode") === null &&
          window.matchMedia("(prefers-color-scheme: dark)").matches)
      : false,

  toggleDarkMode: () =>
    set((state) => {
      const newMode = !state.darkMode;
      if (typeof window !== "undefined") {
        localStorage.setItem("darkMode", newMode.toString());
        if (newMode) {
          document.documentElement.classList.add("dark");
        } else {
          document.documentElement.classList.remove("dark");
        }
      }
      return { darkMode: newMode };
    }),

  initializeTheme: () =>
    set((state) => {
      if (typeof window !== "undefined") {
        const darkModeStored = localStorage.getItem("darkMode");
        const prefersDark = window.matchMedia(
          "(prefers-color-scheme: dark)"
        ).matches;
        const isDarkMode =
          darkModeStored !== null ? darkModeStored === "true" : prefersDark;

        if (isDarkMode) {
          document.documentElement.classList.add("dark");
        } else {
          document.documentElement.classList.remove("dark");
        }
        return { darkMode: isDarkMode };
      }
      return state;
    }),
}));

export const useInitializeTheme = () => {
  const { initializeTheme } = useThemeStore();

  useEffect(() => {
    initializeTheme();
  }, []);
};
