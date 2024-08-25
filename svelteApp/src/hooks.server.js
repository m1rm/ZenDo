export async function handleError({ error, event, status, message }) {
  console.log('Unexpected server side error: ', error)
  return {
    message
  };
}