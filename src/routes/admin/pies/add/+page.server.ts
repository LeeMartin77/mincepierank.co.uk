import type { PageServerLoad } from './$types';
import { getConfig } from '$lib/storage/config';
import { fail, type Actions } from '@sveltejs/kit';
import { validateAdmin } from '../../utilities';
import { writeFile, mkdir } from 'node:fs/promises';
import { setMakerPie, type MakerPie } from '$lib/storage';

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
    const textFields = [
      'year',
      'makerid',
      'id',
      'displayname',
      'fresh',
      'labels',
      'web_link',
      'pack_count', 
      'pack_price_in_pence'
    ];

    const pulledData: { [key:string]: any } = textFields.reduce((acc, curr) => {
      acc[curr] = data.get(curr);
      return acc;
    }, {} as {[key: string]: any})

    const image = data.get('image') as any;

    if (!textFields.every(f => !!pulledData[f]) || !image || !image.stream) {
      return fail(400, { incorrect: true });
    }
    const imgpath = `${pulledData.year}/${pulledData.makerid}/${pulledData.id}/${image.name}`
    const route = `${process.env.IMGPRSSR_DIR || '/imgprssr/'}${imgpath}`;
    await mkdir(route, {
      recursive: true
    });
    await writeFile(`${route}${image.name}`, image.stream());

    pulledData.image_file = imgpath;

    pulledData.labels = pulledData.labels.split(',')

    setMakerPie(pulledData as MakerPie)

    return { success: true };
  }
} satisfies Actions;
