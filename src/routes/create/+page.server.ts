import { error, redirect, type Actions, fail } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getConfig } from '$lib/storage/config';
import { mkdir, writeFile } from 'node:fs/promises';
import { randomUUID } from 'node:crypto';
import { setUserPie, type UserPie } from '$lib/storage';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const load = async (_event: PageServerLoadEvent) => {
  const config = await getConfig();
  if (config.readonly === 'true' || config.customPiesEnabled !== 'true') {
    throw error(404, 'Not Found');
  }
  const activeYear = parseInt(config.activeYear);
  return {
    activeYear
  };
};

export const actions = {
  upload: async ({ request, locals }) => {
    const session = await locals.getSession();
    if (!session?.user?.email) {
      throw error(401, 'Must be signed in!');
    }
    const config = await getConfig();
    if (config.readonly === 'true') {
      throw error(402, 'Currently in readonly mode');
    }

    const data = await request.formData();
    const textFields = [
      'maker',
      'location',
      'displayname',
      'labels',
      'web_link',
      'pack_count',
      'pack_price_in_pence'
    ];

    const pulledData: { [key: string]: any } = textFields.reduce(
      (acc, curr) => {
        acc[curr] = data.get(curr);
        return acc;
      },
      { fresh: false, validated: false } as { [key: string]: any }
    );

    pulledData.fresh = data.get('fresh') == 'on';
    pulledData.clean = false;

    pulledData.year = config.activeYear;
    pulledData.owner_userid = session.user.email;
    pulledData.id = randomUUID();

    const image = data.get('image') as any;

    if (!pulledData.maker || !pulledData.displayname) {
      return fail(400, { incorrect: true });
    }

    if (
      image &&
      image.stream &&
      image.name &&
      image.name !== undefined &&
      image.name !== 'undefined'
    ) {
      const imgpath = `/${pulledData.year}/user-pies/${pulledData.id}`;
      const route = `${process.env.IMGPRSSR_DIR || '/imgprssr'}${imgpath}`;
      await mkdir(route, {
        recursive: true
      });
      await writeFile(`${route}/${image.name}`, image.stream());

      pulledData.image_file = imgpath + '/' + image.name;
    }

    pulledData.labels = pulledData.labels.split(',');

    await setUserPie(pulledData as UserPie);

    throw redirect(303, `/years/${pulledData.year}/user-pies/${pulledData.id}`);
  }
} satisfies Actions;
