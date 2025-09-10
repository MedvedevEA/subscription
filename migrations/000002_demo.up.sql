INSERT INTO public.services("name") 
VALUES ('СберПрайм'),('Яндекс Плюс'),('МТС Premium'),('Т2 Mixx'),('Ozon Premium');
INSERT INTO public.subscriptions(service_id,price,user_id,start_date,stop_date) 
VALUES 
(1,100,'e9c1bc0c-9e9c-413a-84cd-287576e71b25','01.01.2024','01.01.2025'),
(2,200,'e9c1bc0c-9e9c-413a-84cd-287576e71b25','01.06.2024','01.01.2025'),
(3,100,'e9c1bc0c-9e9c-413a-84cd-287576e71b25','01.01.2024','01.01.2024'),
(4,300,'932e9de5-112c-4485-b0cf-0ad0d4cd84db','01.06.2024','01.06.2025'),
(5,400,'006c8b4b-70d4-46e6-8f1f-f53c7c538aa1','01.01.2024',null),
(1,200,'932e9de5-112c-4485-b0cf-0ad0d4cd84db','01.06.2025',null),
(2,500,'006c8b4b-70d4-46e6-8f1f-f53c7c538aa1','01.01.2025','01.01.2025'),
(3,300,'932e9de5-112c-4485-b0cf-0ad0d4cd84db','01.06.2024',null),
(1,100,'d69d5498-0b25-48e0-a3d0-181ea12d6292','01.01.2024','01.01.2025');