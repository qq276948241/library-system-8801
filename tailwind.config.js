/** @type {import('tailwindcss').Config} */

export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,vue}"],
  theme: {
    container: {
      center: true,
      padding: "1.5rem",
    },
    extend: {
      colors: {
        ink: {
          DEFAULT: "#1F3D2B",
          50: "#EEF3EF",
          100: "#D6E2D9",
          200: "#AFC4B6",
          300: "#7E9C88",
          400: "#4D7359",
          500: "#2A5140",
          600: "#1F3D2B",
          700: "#173020",
          800: "#0F1F15",
          900: "#08130B",
        },
        paper: {
          DEFAULT: "#F6F1E6",
          50: "#FBF8F1",
          100: "#F6F1E6",
          200: "#EDE5D3",
          300: "#E0D5BC",
        },
        brass: {
          DEFAULT: "#A8743A",
          300: "#D2A067",
          400: "#C28A4D",
          500: "#A8743A",
          600: "#8A5E2C",
          700: "#6E4920",
        },
        charcoal: {
          DEFAULT: "#2A2620",
          soft: "#5C5347",
          muted: "#8A8073",
        },
      },
      fontFamily: {
        serif: ['"Noto Serif SC"', "Georgia", "serif"],
        sans: ['"Noto Sans SC"', "-apple-system", "BlinkMacSystemFont", "sans-serif"],
        mono: ['"JetBrains Mono"', "ui-monospace", "monospace"],
      },
      boxShadow: {
        card: "0 1px 2px rgba(42,38,32,0.04), 0 8px 24px -12px rgba(31,61,43,0.18)",
        lift: "0 2px 6px rgba(42,38,32,0.06), 0 18px 40px -18px rgba(31,61,43,0.28)",
        inset: "inset 0 1px 0 rgba(255,255,255,0.5)",
      },
      borderRadius: {
        xl2: "1.25rem",
      },
      keyframes: {
        "fade-up": {
          "0%": { opacity: "0", transform: "translateY(10px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        "fade-in": {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
      },
      animation: {
        "fade-up": "fade-up 0.5s cubic-bezier(0.22,1,0.36,1) both",
        "fade-in": "fade-in 0.4s ease both",
      },
    },
  },
  plugins: [],
};
