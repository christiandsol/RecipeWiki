import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
    return { id: Number(params.id) };
};

export const prerender = false;



