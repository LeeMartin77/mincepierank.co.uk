import type { PageServerLoad } from './$types';
import { getConfig } from '$lib/storage/config';
import { fail, type Actions } from '@sveltejs/kit';
import { validateAdmin } from '../utilities';
import { writeFile, mkdir } from 'node:fs/promises';

export const load: PageServerLoad = async () => {
  const config = await getConfig();
  return {
    config
  };
};

export const actions = {
  upload: async ({ request, locals }) => {
    const session = await locals.getSession();
    await validateAdmin(session?.user?.email);
    const data = await request.formData();
    const path = data.get('path') as string;
    const image = data.get('image') as any;
    if (!path || !image || !image.stream) {
      //also validate image is file
      return fail(400, { path, image, incorrect: true });
    }

    const route = `${process.env.IMGPRSSR_DIR || '/imgprssr/'}${path}/`;
    await mkdir(route, {
      recursive: true
    });
    await writeFile(`${route}${image.name}`, image.stream());

    return { success: true };
  }
} satisfies Actions;
