export interface ReportInfo {
  hotelName: string;
  locationName: string;
  task: string;

  text: string;
  images: {
    id: string;
    url: string;
  }[];

  status: string;
}

export const loadReport = async (id: string): Promise<ReportInfo> => {
  return {
    hotelName: "Superhotel",
    locationName: "Москва",
    task: "сделать траляля\nсфотографировать труляля",

    images: [
      {
        id: "1",
        url: "https://www.gentinghotel.co.uk/_next/image?url=https%3A%2F%2Fs3.eu-west-2.amazonaws.com%2Fstaticgh.gentinghotel.co.uk%2Fuploads%2Fhero%2FSuiteNov2022_0008_1920.jpg&w=3840&q=75",
      },
      {
        id: "3",
        url: "https://www.potawatomi.com/application/files/3517/4560/6138/Signature-2-Queen_body.webp",
      },
      {
        id: "4",
        url: "https://images.squarespace-cdn.com/content/v1/60a23657a6164d69e38ddad0/34e78be0-e3f9-4483-9a50-750b62a7b746/PM_FOOD_WEB_3.jpg",
      },
      {
        id: "5",
        url: "https://media.istockphoto.com/id/625006196/photo/sunrise-on-a-tropical-island-palm-trees-on-sandy-beach.jpg?s=612x612&w=0&k=20&c=qGNG4XX4d3SNPDgLgM0GpdEtcPhyldWzQTd38KoC1X8=",
      },
    ],
    text: "Все понравилось\n\nОтличное место\nПравдивые отзывы\nВкусная еда",

    status: "created",
  };
};
