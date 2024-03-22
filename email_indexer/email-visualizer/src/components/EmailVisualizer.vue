<template>
  <div class="flex">
    <div class="w-1/2 bg-gray-100">
      <table class="w-full">
        <thead>
          <tr>
            <th class="px-4 py-2">Subject</th>
            <th class="px-4 py-2">From</th>
            <th class="px-4 py-2">To</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="email in emailList" :key="email.id" @click="showEmailDetails(email)">
            <td class="border px-4 py-2">{{ email.subject }}</td>
            <td class="border px-4 py-2">{{ email.from }}</td>
            <td class="border px-4 py-2">{{ email.to }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="w-1/2 bg-gray-200">
      <div class="p-6" v-if="selectedEmail">
        <h2 class="text-xl font-semibold mb-4">Selected Email</h2>
        <div class="border border-gray-300 p-4 rounded-lg">
          <p><strong>Subject:</strong> {{ selectedEmail.subject }}</p>
          <p><strong>From:</strong> {{ selectedEmail.from }}</p>
          <p><strong>To:</strong> {{ selectedEmail.to }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  // Your script logic here
}
</script>
