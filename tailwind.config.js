module.exports = {
  mode: 'jit',
  content: ['./public/**/*.html'],
  purge: ['./public/**/*.html'],
  theme: {
    extend: {
      fontFamily: {
        customFont: ['"Custom Font"', "raleway"],
        // Add more custom font families as needed
      },
    },
  },
  plugins: [],
}
