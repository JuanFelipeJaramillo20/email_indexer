<template>
    <NavbarDark />
    <div id="about" class="w-100 min-w-screen min-h-screen bg-gray-900 text-white">
      <!-- Navbar -->
  
      <!-- Main content -->
      <div class="container mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold mb-4">About Us</h1>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Project Description -->
          <div>
            <h2 class="text-xl font-semibold mb-2">Project Description</h2>
            <p class="text-gray-300 leading-relaxed">
              Your project description goes here. Explain the purpose, goals, and key features of your web app.
            </p>
          </div>
          <!-- Technologies Used -->
          <div>
            <h2 class="text-xl font-semibold mb-2">Technologies Used</h2>
            <ul class="list-disc list-inside text-gray-300">
              <li>Tailwind CSS</li>
              <li>Vue.js</li>
              <li>...</li>
            </ul>
          </div>
        </div>
        <!-- Additional Information -->
        <div class="mt-8">
          <h2 class="text-xl font-semibold mb-2">Additional Information</h2>
          <p class="text-gray-300 leading-relaxed">
            Add any additional information about your project here. This could include links, credits, acknowledgments, etc.
          </p>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import NavbarDark from "../components/NavbarDark.vue";
  
  export default {
    components: {
      NavbarDark
    },
    data() {
      return {};
    }
  };
  </script>
  