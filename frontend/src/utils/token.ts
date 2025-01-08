// utils/token.ts
export const setAuthToken = (token: string) => {
  console.log("auth Token: ", token)
  localStorage.setItem('authToken', token)
}

export const getAuthToken = () => {
  return localStorage.getItem('authToken')
}

export const removeAuthToken = () => {
  localStorage.removeItem('authToken')
}