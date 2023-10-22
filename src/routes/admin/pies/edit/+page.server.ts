import type { PageServerLoad } from './$types';
import { getConfig } from '$lib/storage/config';
import { fail, type Actions, redirect } from '@sveltejs/kit';
import { validateAdmin } from '../../utilities';
import { writeFile, mkdir } from 'node:fs/promises';
import {
  setMakerPie,
  type MakerPie,
  getPieByMakerAndId,
  deletePieByMakerAndId,
  getMincePieMakers
} from '$lib/storage';

export const load: PageServerLoad = async (event) => {
  const config = await getConfig();
  const year = parseInt(event.url.searchParams.get('year') || '0');
  const makerid = event.url.searchParams.get('makerid');
  const id = event.url.searchParams.get('id');

  const editing = !!year && year && !!makerid && !!id;

  let pie;

  if (editing) {
    pie = (await getPieByMakerAndId(year, makerid, id)).unwrapOr(undefined);
  }

  const makers = (await getMincePieMakers()).unwrapOr([]);

  return {
    config,
    editing,
    pie,
    makers
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
    pulledData.validated = data.get('validated') == 'on';

    const image = data.get('image') as any;

    if (!textFields.every((f) => !!pulledData[f]) || !image || !image.stream) {
      return fail(400, { incorrect: true });
    }
    const imgpath = `/${pulledData.year}/${pulledData.makerid}/${pulledData.id}`;
    const route = `${process.env.IMGPRSSR_DIR || '/imgprssr'}${imgpath}`;
    await mkdir(route, {
      recursive: true
    });
    await writeFile(`${route}/${image.name}`, image.stream());

    pulledData.image_file = imgpath + '/' + image.name;

    pulledData.labels = pulledData.labels.split(',');

    await setMakerPie(pulledData as MakerPie);

    return { success: true };
  },
  edit: async ({ request, locals }) => {
    const session = await locals.getSession();
    await validateAdmin(session?.user?.email);
    const data = await request.formData();
    const pieData = await getPieByMakerAndId(
      parseInt((data.get('year') as string) || ''),
      data.get('makerid') as string,
      data.get('id') as string
    );

    if (pieData.isErr()) {
      console.log(data);
      console.log(data.get('year'));
      console.log(data.get('makerid'));
      console.log(data.get('id'));
      console.error('Err loading pie');
      return fail(400, { incorrect: true });
    }
    const pie = pieData.value;

    const displayname = data.get('displayname');
    if (displayname) {
      pie.displayname = displayname as string;
    }
    const web_link = data.get('web_link');
    if (web_link) {
      pie.web_link = web_link as string;
    }
    const pack_count = parseInt((data.get('pack_count') as string) || '');
    if (pack_count) {
      pie.pack_count = pack_count;
    }
    const pack_price_in_pence = parseInt((data.get('pack_price_in_pence') as string) || '');
    if (pack_price_in_pence) {
      pie.pack_price_in_pence = pack_price_in_pence;
    }

    pie.fresh = (data.get('fresh') as any) === 'on';

    pie.validated = (data.get('validated') as any) === 'on';

    const image = data.get('image') as any;

    if (
      image &&
      image.stream &&
      image.name &&
      image.name !== undefined &&
      image.name !== 'undefined'
    ) {
      const imgpath = `/${pie.year}/${pie.makerid}/${pie.id}`;
      const route = `${process.env.IMGPRSSR_DIR || '/imgprssr'}${imgpath}`;
      await mkdir(route, {
        recursive: true
      });
      await writeFile(`${route}/${image.name}`, image.stream());

      pie.image_file = imgpath + '/' + image.name;
    }

    pie.labels = ((data.get('labels') as string) || '').split(',');

    await setMakerPie(pie as MakerPie);

    throw redirect(303, `/admin/pies/edit?year=${+pie.year}&makerid=${pie.makerid}&id=${pie.id}`);
  },
  delete: async ({ request, locals }) => {
    const session = await locals.getSession();
    await validateAdmin(session?.user?.email);
    const data = await request.formData();
    const pieData = await getPieByMakerAndId(
      parseInt((data.get('year') as string) || ''),
      data.get('makerid') as string,
      data.get('id') as string
    );

    if (pieData.isErr()) {
      console.log(data);
      console.log(data.get('year'));
      console.log(data.get('makerid'));
      console.log(data.get('id'));
      console.error('Err loading pie');
      return fail(400, { incorrect: true });
    }

    await deletePieByMakerAndId(
      parseInt((data.get('year') as string) || ''),
      data.get('makerid') as string,
      data.get('id') as string
    );

    throw redirect(303, `/admin/pies/`);
  }
} satisfies Actions;
