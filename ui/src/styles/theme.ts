import { extendTheme } from '@chakra-ui/react'

export const colors = {
  1: '#FC2D81', //pink
  2: '#0F1F48', //dark
  3: 'rgba(15, 31, 72, 0.3)', //gray
  4: '#F9F9F9', //light gray
  5: '#FFFFFF',
  6: '#F2297B', //pink hover
  7: '#386FFE', //blue
  8: '#2F66F1', //blue hover
}

const breakpoints = {
  sm: '375px',
  md: '768px',
  lg: '1024px',
  lg1: '1280px',
  xl: '1920px',
  '2xl': '2560px',
}

const fonts = {
  body: 'Roboto Slab, serif',
  heading: 'Roboto Slab, serif',
}

export const components = {
  Header: {
    baseStyle: () => ({
      w: '100%',
      h: '98px',
      justifyContent: 'flex-start',
      alignItems: 'center',
      fontSize: '1.875rem',
      gap: 3,
    }),
  },
  Button: {
    baseStyle: {
      w: '260px',
      h: '50px',
      borderRadius: '24px',
      color: colors[5],
      fontWeight: 400,
      fontSize: '18px',
      lineHeight: '24px',
      _hover: {
        _disabled: {
          backgroundColor: colors[4],
          color: colors[3],
          opacity: 1,
        },
      },
      _disabled: {
        backgroundColor: colors[4],
        color: colors[3],
        opacity: 1,
      },
    },
    variants: {
      primary: {
        bg: colors[1],
        backgroundColor: colors[1],
        _hover: {
          bg: colors[6],
          backgroundColor: colors[6],
        },
      },
      secondary: {
        bg: colors[7],
        backgroundColor: colors[7],
        _hover: {
          bg: colors[8],
          backgroundColor: colors[8],
        },
      },
    },
  },
}

const theme = extendTheme({
  colors,
  components,
  breakpoints,
  fonts,
})

export { theme }
