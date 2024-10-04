/**
 * @param {{ error: any, event: any, status: number, message: string }} params
 */
export async function handleError({ error, event, status, message }) {
  console.log('Unexpected server side error: ', error);
  return {
    message
  };
}
