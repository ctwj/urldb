/**
 * Sample Plugin Main Entry Point
 * This is the main entry point for the sample plugin
 */

// Plugin initialization
console.log('Sample plugin loaded successfully!');

// Plugin can export functions or objects
module.exports = {
  name: 'sample-plugin',
  version: '1.0.0',

  // Plugin initialization function
  init() {
    console.log('Sample plugin initialized');
  },

  // Plugin cleanup function
  cleanup() {
    console.log('Sample plugin cleaned up');
  }
};