import { redirect } from 'next/navigation';

export default function Home() {
  // Redirect to the default seeded product ID to demonstrate the flow
  redirect('/item/MLA43960787');
}
