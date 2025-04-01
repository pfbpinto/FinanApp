
/* SEEDER INFORMATION: The tables who are being feed by a seeeder shoulnt change the structure in the future, there ones without a seeder might suffer significant changes, as they haven't been refined and executed a Proof fo Concept (POC)
 on them
 
 Seeder User:pfbpinto@hotmail.com
 Password:123
 
 
 */


 -- Entity Seder

INSERT INTO Entity (EntityName, EntityType,EntityCategory,IsActive) VALUES
/*1*/('Group','Group Income','Parent',TRUE ),
/*2*/('Group','Group Expense','Parent',TRUE),
/*3*/('Group Income','Group Income Tax','Child',TRUE),
/*4*/('Group Income','Group Income Expense','Child',TRUE),
/*5*/('User','User Income','Parent',TRUE),
/*6*/('User','User Expense','Parent',TRUE),
/*7*/('User Income','User Income Tax','Child',TRUE),
/*8*/('User Income','User Income Expense','Child',TRUE),
/*9*/('Asset','Asset Tax','Parent',TRUE),
/*10*/('Asset','Asset Expense','Parent',TRUE),
/*11*/('Asset','Asset Income','Parent',TRUE),
/*12*/('Asset Income','Asset Income Tax','Child',TRUE),
/*13*/('Asset Income','Asset Income Expense','Child',TRUE);



INSERT INTO IncomeType (IncomeTypeName, IncomeDescription, EntityID) VALUES
/*1*/('Salary', 'Income provenient of a service for a company or person, usually recurring for a fixed amount.',1),
/*2*/('Rent', 'Income from an owning asset (vehicle, real state…) for a predefined amount of time',11),
/*3*/('Provided Service', 'Income from an one time service, receiving the amount agreed between both parties',1),
/*4*/('Investment', 'Income from any type of investment, through stock, bank, crypto, or any associated activities which returns an amount of money, based on an initial funding.',1),
/*5*/('Loan', 'Income incoming from an entity that provided funding, expecting the return of the amount, usually with interest.',1),
/*6*/('Dividends', 'Income provenient of a company.',1);

-- Tax Type Seder
INSERT INTO TaxType (TaxTypeName, TaxDescription, TaxCountry, TaxJurisdiction, TaxPercentage,EntityID) 
VALUES
/*1*/('II', 'Import tax for goods coming from outside the country.', 'Brazil', 'Federal', NULL,1),
/*2*/('IOF', 'Tax on financial transactions, for loans, shares and other financial actions', 'Brazil', 'Federal', NULL,1),
/*3*/('IRPF', 'Individual Income Tax, on the citizens income', 'Brazil', 'Federal', NULL,10),
/*4*/('COFINS', 'Social security financing contribution', 'Brazil', 'Federal', NULL,1),
/*5*/('PIS', 'Social Integration Program', 'Brazil', 'Federal', NULL,1),
/*6*/('CSLL', 'Social contribution on net profit', 'Brazil', 'Federal', NULL,1),
/*7*/('INSS', 'National Institute of Social Security', 'Brazil', 'Federal', NULL,1),
/*8*/('IRPJ', 'Corporate Income Tax, on the income of CNPJs', 'Brazil', 'Federal', NULL,1),
/*9*/('ICMS', 'Taxes on circulation of goods and services', 'Brazil', 'State', NULL,1),
/*10*/('IPVA', 'Tax on motor ownership automotive', 'Brazil', 'State', NULL,1),
/*11*/('ITCMD', 'Tax on inheritance and donation', 'Brazil', 'State', NULL,1),
/*12*/('IPTU', 'Tax on urban land property', 'Brazil', 'City', NULL,12),
/*13*/('ISS', 'Tax on services', 'Brazil', 'City', NULL,1),
/*15*/('ITBI', 'Tax on transfer of real estate', 'Brazil', 'City', NULL,1);


-- Expense Type Seder
INSERT INTO ExpenseType (ExpenseTypeName, ExpenseDescription,EntityID) VALUES
/*1*/('Bill', 'Regular utility payments such as water, electricity, internet services, phone lines, and domestic services like cleaning or maintenance.',1),
/*2*/('Mortgage', 'Long-term financial commitments related to the purchase of property, including homes, apartments, buildings, vehicles, or large appliances and furniture.',1),
/*3*/('Rent', 'House, apartment, building rent',1),
/*4*/('Entertainment', 'Costs associated with leisure and recreation, including dining out, entertainment activities, vacations, travel tours, and other forms of enjoyment.',1),
/*5*/('Food and Supply', 'Purchases for everyday household necessities, including groceries, household supplies, and personal care items.',1),
/*6*/('Exchange Costs','Costs on currency exchange',11);


-- Asset Type Seder

INSERT INTO AssetType (AssetTypeName, AssetTypeDescription, EntityID) VALUES
/*1*/('Real Estate', 'Residential property (houses, apartments, condos, vacation homes, etc.), Commercial property (office buildings, retail space, etc.), Land or undeveloped property, Timeshares, Rental properties, Vacation properties',1),
/*2*/('Investment', 'Stocks, bonds, and securities, Mutual funds, ETFs, index funds, Cryptocurrencies (Bitcoin, Ethereum, etc.), Commodities (gold, silver, oil, etc.), Peer-to-peer lending, Retirement accounts (pension funds), Cash savings (savings accounts, CDs, money market funds), Private equity or venture capital investments, REITs (Real Estate Investment Trusts)',1),
/*3*/('Business Ownership', 'Ownership in companies (shares, equity, LLC ownership), Franchise interests, Sole proprietorship or partnership business value, Intellectual Property (patents, trademarks, copyrights)',1),
/*4*/('Vehicle', 'Vehicles (cars, boats, motorcycles, RVs)',1);

-- Currency Seder

INSERT INTO Currency ( CurrencyName,CurrencyAbreviation, CurrencySymbol) VALUES
('Brazilian Real','BRL','R$'); 

-- User Seeder
INSERT INTO Recurrency (RecurrencyName, RecurrencyPeriod) VALUES 
/*1*/('One Time', 'today'),
/*2*/('Monthly', '1 month'),
/*3*/('Quarterly','4 months'),
/*4*/('Yearly', '1 year')
/*5*/('Variable','undefined');
/*-- ATUALIZANDO A TABELA RECURRENCY PRA FUTURAMENTE UPGRADE A PROCEDURE PARA UTILIZAR O PERIDO DIRETAMENTE DA TABELA COMO EXEMPLO ABAIXO:
      -- Configura número de inserções e intervalo conforme a recorrência
        IF p_RecurrencyID = 1 THEN -- One time
            v_Iterations := 1;
            v_Increment := '1 day'::INTERVAL;
        ELSIF p_RecurrencyID = 2 THEN -- Monthly
            v_Iterations := 12;
            v_Increment := '1 month'::INTERVAL;
        ELSIF p_RecurrencyID = 3 THEN -- Quarterly
            v_Iterations := 4;
            v_Increment := '4 months'::INTERVAL;
        ELSIF p_RecurrencyID = 4 THEN -- Yearly
            v_Iterations := 1;
            v_Increment := '1 year'::INTERVAL;
        ELSE
            RAISE EXCEPTION 'Invalid RecurrencyID';
        END IF;

        -- Insere múltiplos registros conforme a recorrência
        v_CurrentDate := p_BeginDate;
        FOR i IN 1..v_Iterations LOOP
            -- Define a data do próximo registro para calcular a end date
            v_NextDate := v_CurrentDate + v_Increment;
 */           
-- User Seeder
INSERT INTO UserProfile (FirstName, LastName, DateOfBirth, UserPassword, EmailAddress, UserSubscription) 
VALUES 
('Pedro', 'Pinto', '1985-07-30', '$2a$10$"$2a$10$gFxlAuEUcogaJS4R8jiOiudfVjM3q0H.w9GsLvHPuoIcqKtoH2vE6', 'pfbpinto@hotmail.com', '0');


--User AssetSeeder
INSERT INTO UserAsset (AssetTypeID, UserProfileID, UserAssetName, UserAssetValueAmount, UserAssetAcquisitionBeginDate, UserAssetAcquisitionEndDate, IsActive) VALUES 
(1, 1, 'Apartamento Onix', 100000.00, '2020-01-01', NULL, TRUE);


--User Category Seeder
INSERT INTO UserCategory (UserCategoryName, UserProfileID, EntityID, IsActive) VALUES 
/*1*/('Recebimentos', 1, 5,true),
/*2*/('Compras Gerais', 1, 6,true),
/*3*/('Impostos Pessoais', 1,7,true),
/*4*/('Gastos Obrigatorios', 1,8,true),
/*5*/('Impostos Imobiliarios', 1,12,true),
/*6*/('Financiamentos', 1, 13,true),
/*7*/('Fontes de Renda Imobiliaria', 1, 11,true);


-- FinancialUserItemName Seeder
INSERT INTO  FinancialUserItem (FinancialUserItemName, EntityID,UserEntityID, RecurrencyID, FinancialUserEntityItemID,ParentFinancialUserItemID,IsActive) VALUES 
/*1*/('Salario Exterior', 5,1,1,2,NULL,TRUE),
/*2*/('Novo Laptop', 6,1,1,1,NULL,TRUE),
/*3*/('Imposto de Renda Salario', 7,1,1,2,1,TRUE),
/*4*/('Cambio do Salario', 8,1,2,6,1,TRUE),
/*5*/('IPTU Apartamento Onix',10,1,4,12,NULL,TRUE),
/*6*/('Financiamento Apartamento Onix',10,1,2,2,NULL,TRUE),
/*7*/('Aluguel Apartamento Onix',11,1,2, 2,NULL,TRUE);

--Forecast Seeder
INSERT INTO UserFinancialForecast (UserCategoryID, FinancialUserItemID, UserFinancialForecastAmount,UserFinancialForecastBeginDate, UserFinancialForecastEndDate, CurrencyID) 
VALUES 
(1, 1,8500, '2025-02-01', NULL, 1),
(2, 2,3500, '2025-03-05', '2025-01-05', 1),
(3, 3,250, '2025-03-01', NULL, 1),
(4, 4,50, '2025-03-01', NULL, 1),
(5, 5,1850, '2025-03-30', NULL, 1),
(6, 6,1120, '2025-03-15', NULL, 1),
(7, 7,1900, '2025-03-10', NULL, 1),
(1, 1,8500, '2025-04-01', NULL, 1),
(3, 3,250, '2025-04-01', NULL, 1),
(4, 4,50, '2025-04-01', NULL, 1),
(5, 5,1850, '2025-04-30', NULL, 1),
(6, 6,1120, '2025-04-15', NULL, 1),
(7, 7,1900, '2025-04-10', NULL, 1);

INSERT INTO UserFinancialActual (UserCategoryID, FinancialUserItemID,UserFinancialActualAmount,UserFinancialActualtBeginDate, UserFinancialActualEndDate, CurrencyID, Note) VALUES 
(1, 1, 8500 , '2025-01-01' , NULL		  , 1,NULL),
(2, 2, 3500 , '2025-01-05' , '2025-01-05', 1,NULL),
(3, 3, 234  , '2025-01-01' , NULL		  , 1,NULL),
(4, 4, 80   , '2025-01-01' , NULL		  , 1,'Dolar em alta'),
(5, 5, 1450 , '2025-01-30' , NULL		  , 1,'Desconto por Antecipação'),
(6, 6, 1120 , '2025-01-15' , NULL	      , 1,NULL),
(7, 7, 2100 , '2025-01-20' , NULL        , 1,'Atraso de Pagamento Inquilino');

INSERT INTO UserForecastActualRelation (UserFinancialActualID, UserFinancialForecastID) 
VALUES 
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5),
(6, 6),
(7, 7);