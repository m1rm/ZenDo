export async function handleError({ error, event, status, message }) {
  console.log('Unexpected client side error: ', error)
  return {
    message
};
}