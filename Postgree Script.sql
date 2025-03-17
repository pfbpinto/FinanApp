
/*postgree version*/



----------------------------------------------
-------------- REFERENCES -----------------------
-------------------------------------------------

-- User Profile: Named as user profile as "User" is a native command in sql, to avoid issues, changed to userProfile
CREATE TABLE UserProfile (
    UserProfileID SERIAL PRIMARY KEY,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(255) NOT NULL,
	DateOfBirth DATE NOT NULL,
    UserPassword VARCHAR(150) NOT NULL,
    EmailAddress VARCHAR(255) UNIQUE NOT NULL,
    UserSubscription BOOLEAN NOT NULL, -- Boolean Active/inactive
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Income Type
CREATE TABLE IncomeType (
    IncomeTypeID SERIAL PRIMARY KEY,
    IncomeTypeName VARCHAR(100) NOT NULL,
    IncomeDescription VARCHAR(255) NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO IncomeType (IncomeTypeName, IncomeDescription) VALUES
('Salary', 'Income provenient of a service for a company or person, usually recurring for a fixed amount.'),
('Rent', 'Income from an owning asset (vehicle, real state…) for a predefined amount of time'),
('Provided Service', 'Income from an one time service, receiving the amount agreed between both parties'),
('Investment', 'Income from any type of investment, through stock, bank, crypto, or any associated activities which returns an amount of money, based on an initial funding.'),
('Loan', 'Income incoming from an entity that provided funding, expecting the return of the amount, usually with interest.'),
('Dividends', 'Income provenient of a company.');


-- Tax Type
CREATE TABLE TaxType (
    TaxTypeID SERIAL PRIMARY KEY,
    TaxTypeName VARCHAR(100) NOT NULL,
    TaxDescription VARCHAR(255) NOT NULL,
    TaxCountry VARCHAR(100) NOT NULL,
    TaxJurisdiction VARCHAR(100) NOT NULL,
    TaxPercentage DECIMAL(5,2),
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO TaxType (TaxTypeName, TaxDescription, TaxCountry, TaxJurisdiction, TaxPercentage) 
VALUES
('II', 'Import tax for goods coming from outside the country.', 'Brazil', 'Federal', NULL),
('IOF', 'Tax on financial transactions, for loans, shares and other financial actions', 'Brazil', 'Federal', NULL),
('IRPF', 'Individual Income Tax, on the citizens income', 'Brazil', 'Federal', NULL),
('COFINS', 'Social security financing contribution', 'Brazil', 'Federal', NULL),
('PIS', 'Social Integration Program', 'Brazil', 'Federal', NULL),
('CSLL', 'Social contribution on net profit', 'Brazil', 'Federal', NULL),
('INSS', 'National Institute of Social Security', 'Brazil', 'Federal', NULL),
('IRPJ', 'Corporate Income Tax, on the income of CNPJs', 'Brazil', 'Federal', NULL),
('ICMS', 'Taxes on circulation of goods and services', 'Brazil', 'State', NULL),
('IPVA', 'Tax on motor ownership automotive', 'Brazil', 'State', NULL),
('ITCMD', 'Tax on inheritance and donation', 'Brazil', 'State', NULL),
('IPTU', 'Tax on urban land property', 'Brazil', 'City', NULL),
('ISS', 'Tax on services', 'Brazil', 'City', NULL),
('ITBI', 'Tax on transfer of real estate', 'Brazil', 'City', NULL);

-- Expense Type
CREATE TABLE ExpenseType (
    ExpenseTypeID SERIAL PRIMARY KEY,
    ExpenseTypeName VARCHAR(100) NOT NULL,
    ExpenseDescription VARCHAR(255) NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO ExpenseType (ExpenseTypeName, ExpenseDescription) VALUES
('Bill', 'Regular utility payments such as water, electricity, internet services, phone lines, and domestic services like cleaning or maintenance.'),
('Mortgage', 'Long-term financial commitments related to the purchase of property, including homes, apartments, buildings, vehicles, or large appliances and furniture.'),
('Rent', 'House, apartment, building rent'),
('Entertainment', 'Costs associated with leisure and recreation, including dining out, entertainment activities, vacations, travel tours, and other forms of enjoyment.'),
('Food and Supply', 'Purchases for everyday household necessities, including groceries, household supplies, and personal care items.');

-- Asset Type
CREATE TABLE AssetType (
    AssetTypeID SERIAL PRIMARY KEY,
    AssetTypeName VARCHAR(100) NOT NULL,
    AssetTypeDescription VARCHAR(255) NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO AssetType (AssetTypeName, AssetTypeDescription) VALUES
('Real Estate', 'Residential property (houses, apartments, condos, vacation homes, etc.), Commercial property (office buildings, retail space, etc.), Land or undeveloped property, Timeshares, Rental properties, Vacation properties'),
('Investment', 'Stocks, bonds, and securities, Mutual funds, ETFs, index funds, Cryptocurrencies (Bitcoin, Ethereum, etc.), Commodities (gold, silver, oil, etc.), Peer-to-peer lending, Retirement accounts (pension funds), Cash savings (savings accounts, CDs, money market funds), Private equity or venture capital investments, REITs (Real Estate Investment Trusts)'),
('Business Ownership', 'Ownership in companies (shares, equity, LLC ownership), Franchise interests, Sole proprietorship or partnership business value, Intellectual Property (patents, trademarks, copyrights)'),
('Vehicle', 'Vehicles (cars, boats, motorcycles, RVs)');

-- Currency

CREATE TABLE Currency (
    CurrencyID SERIAL PRIMARY KEY,
    CurrencyName VARCHAR(100) NOT NULL,
	CurrencyAbreviation VARCHAR(10) NOT NULL,
	CurrencySymbol VARCHAR(5) NOT NULL,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


--Exchange Rate always based on Dollar
CREATE TABLE CurrencyExchangeRate (
    CurrencyExchangeRateID SERIAL PRIMARY KEY,
	CurrencyID INT NOT NULL, --FK Currency
	ExchangeRateValue DECIMAL(15,2) NOT NULL, 
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID) ON DELETE CASCADE

);


-------------------------------------------------
-------------- DOCUMENT -----------------------
-------------------------------------------------
-- Document Type Scope
-- Scope covers the entity which this document is part of. It can be from User, Group, or line item (tax, income, expense)

CREATE TABLE DocumentScopeType (
    DocumentScopeTypeID SERIAL PRIMARY KEY,
    DocumentScopeName VARCHAR(100) NOT NULL,
    DocumentScopeDescription VARCHAR(255) NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Document Type
CREATE TABLE DocumentType (
    DocumentTypeID SERIAL PRIMARY KEY,
    DocumentScopeTypeID INT NOT NULL,
    DocumentTypeName VARCHAR(100) NOT NULL,
    DocumentDescription VARCHAR(255) NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_DocumentType_DocumentScopeType FOREIGN KEY (DocumentScopeTypeID) REFERENCES DocumentScopeType(DocumentScopeTypeID)
);

-- Generalized Document Table
CREATE TABLE Document (
    DocumentID SERIAL PRIMARY KEY,
	DocumentReferenceNumber VARCHAR(150), -- For personal documentation or documentats that have specific naming convetion
    DocumentURL TEXT,
    DocumentPath TEXT,
    DocumentTypeID INT NOT NULL,
    EntityTypeName VARCHAR(50) NOT NULL, -- --Open text: "Group", "User", "User Asset", or "Group Asset "
    EntityTypeID INT NOT NULL,-- If "Group"=FinancialGroupID, If "User"=UserProfileID, If "Group Asset"=FinancialGroupAssetID , If "User Asset"=UserAssetID
	EntityItemTypeName VARCHAR(50), -- Open textfield= Tax,Income, Expense
    EntityItemTypeNameID INT,-- TypeID from originating table
    DocumentValidDate DATE NOT NULL,
	DocumentExpirationDate DATE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_Document_DocumentType FOREIGN KEY (DocumentTypeID) REFERENCES DocumentType(DocumentTypeID)
);

-------------------------------------------------
-------------- GROUP MANAGEMENT -----------------
-------------------------------------------------
-- Important: Referenced as "Financial Group" as "Group" is a native key word in SQL

-- Financial Group -- The Group Information and who ownes it (which user owns it)


CREATE TABLE FinancialGroup (
    FinancialGroupID SERIAL PRIMARY KEY,
    FinancialGroupName VARCHAR(255) NOT NULL,
    FinancialGroupOwnerID INT NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_FinancialGroup_UserProfile FOREIGN KEY (FinancialGroupOwnerID) REFERENCES UserProfile(UserProfileID) ON DELETE CASCADE
);

-- Group Category
CREATE TABLE GroupCategory (
    GroupCategoryID SERIAL PRIMARY KEY,
    GroupCategoryName VARCHAR(255) NOT NULL,
    FinancialGroupID INT NOT NULL, -- FK
    ItemTypeName VARCHAR(50) NOT NULL,--- Open text field= Tax,Income, Expense, Asset
    ItemTypeNameID INT NOT NULL,-- TypeID from originating table
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserCategory_FinancialGroup FOREIGN KEY (FinancialGroupID) REFERENCES FinancialGroup(FinancialGroupID) ON DELETE CASCADE
);


-- Financial Group User - who is inside thar group
CREATE TABLE FinancialGroupUser (
    FinancialGroupUserID SERIAL PRIMARY KEY,
    FinancialGroupID INT NOT NULL, -- FK
    UserProfileID INT NOT NULL, --FK
	ActiveUser BOOLEAN NOT NULL, -- You can add non-finanapp users like children, wife, ...
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_FinancialGroupUser_FinancialGroup FOREIGN KEY (FinancialGroupID) REFERENCES FinancialGroup(FinancialGroupID) ON DELETE CASCADE,
    CONSTRAINT FK_FinancialGroupUser_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID)
);

-- Asset owned by a group -- assets owned by the group
CREATE TABLE FinancialGroupAsset (
    FinancialGroupAssetID SERIAL PRIMARY KEY,
	AssetTypeID INT NOT NULL, -- FK
	FinancialGroupID INT NOT NULL, -- FK
    FinancialGroupAssetName VARCHAR(100) NOT NULL, 
    FinancialGroupAssetAmountValue DECIMAL(15,2) NOT NULL,
    FinancialGroupAssetAssetAcquisitionBeginDate DATE NOT NULL,
    FinancialGroupAssetAssetAcquisitionEndDate DATE,
	CONSTRAINT FK_FinancialGroupAsset_FinancialGroup FOREIGN KEY (FinancialGroupID) REFERENCES FinancialGroup(FinancialGroupID) ON DELETE CASCADE,
	CONSTRAINT FK_FinancialGroupAsset_AssetType FOREIGN KEY (AssetTypeID) REFERENCES AssetType(AssetTypeID)
);


-- As asset and income are an entity on it's own, required details around the asset and incomes to define the owned % by User
CREATE TABLE FinancialGroupUserOwnership (
    FinancialGroupUserOwnershipID SERIAL PRIMARY KEY,
    EntityTypeName VARCHAR(50) NOT NULL,       -- Allowed Values: "Group Asset" or "Group Income"
	EntityTypeID INT NOT NULL,                 -- Allowed Values:
											   -- If Entity Name="Group Asset"  FinancialGroupAssetID
									           -- If Entity Name="Group Income" FinantialGroupForecastID 						             
	EntityItemTypeName VARCHAR(50) NOT NULL, -- Allowed Values:
											 -- If Entity Name="Group Asset" then "Income", "Tax" or "Expense"
									         -- If Entity Name="Group Income" then "Tax" or "Expense"
    EntityItemTypeID INT NOT NULL,--If EntityItemTypeName="Group Income" Then IncomeTypeID
							      --If EntityItemTypeName="Expense" Then ExpenseTypeID
							      --If EntityItemTypeName-"Tax" Then TaxTypeID
	UserProfileID INT NOT NULL, -- FK
	FinancialGroupID INT NOT NULL, -- FK
	UserOwningPercentage DECIMAL(5,2) NOT NULL,
	OwnershipStartDate DATE NOT NULL,
	OwnershipEndDate DATE,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_FinancialGroupAssetOwnership_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID),
	CONSTRAINT FK_FinancialGroupAssetOwnership_FinancialGroup FOREIGN KEY (FinancialGroupID) REFERENCES FinancialGroup(FinancialGroupID)
);


-- Financial Group Forecast 
CREATE TABLE FinantialGroupForecast (
    FinantialGroupForecastID SERIAL PRIMARY KEY,
	EntityTypeName VARCHAR(50) NOT NULL, --Open text: "Group", "Group Income" or "Group Asset "
	EntityTypeID INT NOT NULL, -- If "Group"=FinancialGroupID, IF "Group Asset"=FinancialGroupAssetID , If "Group Income"=FinanitalGroupForecastID
	EntityItemTypeName VARCHAR(50) NOT NULL, -- If Entity Name="Group" then "Income or Expense"
									   -- If Entity Name="Group Income" then "Tax or Expense"
									   -- If Entity Name="Group Asset" then "Income, Expense or Tax"
    EntityItemTypeID INT NOT NULL,--If EntityItemTypeName="Income" Then IncomeTypeID
							--If "Expense" Then ExpenseTypeID
							--If "Tax" then TaxTypeID
    FinantialGroupForecastAmount DECIMAL(15,2) NOT NULL, --Amount Value
    FinantialGroupForecastBeginDate DATE NOT NULL,
    FinantialGroupForecastEndDate DATE,
	CurrencyID INT, -- FK
	GroupCategoryID INT, --FK
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_FinantialGroupForecast_CurrencyID FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID),
	CONSTRAINT FK_FinantialGroupForecast_GroupCategory FOREIGN KEY (GroupCategoryID) REFERENCES GroupCategory(GroupCategoryID)
);

-- Financial Group Actuals
CREATE TABLE FinantialGroupActual(
    FinantialGroupActualID SERIAL PRIMARY KEY,
	EntityTypeName VARCHAR(50) NOT NULL, --Open text: "Group", "Group Income" or "Group Asset "
	EntityTypeID INT NOT NULL, -- If "Group"=FinancialGroupID, IF "Group Asset"=FinancialGroupAssetID , If "Group Income"=FinantialGroupActualID
	EntityItemTypeName VARCHAR(50) NOT NULL, -- If Entity Name="Group" then "Income or Expense"
									   -- If Entity Name="Group Income" then "Tax or Expense"
									   -- If Entity Name="Group Asset" then "Income, Expense or Tax"
    EntityItemTypeID INT NOT NULL,--If EntityItemTypeName="Income" Then IncomeTypeID
							--If "Expense" Then ExpenseTypeID
							--If "Tax" then TaxTypeID
    FinantialGroupActualAmount DECIMAL(15,2) NOT NULL, --Amount Value
    FinantialGroupActualBeginDate DATE NOT NULL,
    FinantialGroupActualEndDate DATE,
	CurrencyID INT, -- FK
	FinantialGroupForecastID INT, --FK of the Forecast if exist
	GroupCategoryID INT, -- FK
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_FinancialGroupAsset_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID),
	CONSTRAINT FK_FinancialGroupAsset_Forecast FOREIGN KEY (FinantialGroupForecastID) REFERENCES FinantialGroupForecast(FinantialGroupForecastID),
    CONSTRAINT FK_FinancialGroupAsset_GroupCategory FOREIGN KEY (GroupCategoryID) REFERENCES GroupCategory(GroupCategoryID)
);


-------------------------------------------------
-------------- USER MANAGEMENT -----------------
-------------------------------------------------



-- User Category
CREATE TABLE UserCategory (
    UserCategoryID SERIAL PRIMARY KEY,
    UserCategoryName VARCHAR(255) NOT NULL,
    UserProfileID INT NOT NULL,
    ItemTypeName VARCHAR(50) NOT NULL,--- Open text field= Tax,Income, Expense, Asset
    ItemTypeNameID INT NOT NULL,-- TypeID from originating table
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserCategory_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID) ON DELETE CASCADE
);

INSERT INTO UserCategory (UserCategoryName, UserProfileID, ItemTypeName, ItemTypeNameID)
VALUES
('Tax Payer', 3, 'Tax', 1);  -- 3 é o UserProfileID, 1 é o ID do tipo de item (pode ser de 'TaxType' ou outra tabela dependendo de como você gerencia esses IDs).



-- Asset Management. as Asset is a complex entity on it's own, will be manage separately. On MVP, assets will either be represented on one or other. If it's owned buy the user

CREATE TABLE UserAsset (
    UserAssetID SERIAL PRIMARY KEY,
	AssetTypeID INT NOT NULL, -- FK
	UserProfileID INT NOT NULL, -- FK
    AssetName VARCHAR(100) NOT NULL, 
	AssetValueAmount DECIMAL(15,2) NOT NULL,
    AssetAcquisitionBeginDate DATE NOT NULL,
    AssetAcquisitionEndDate DATE,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_UserAsset_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID) ON DELETE CASCADE,
	CONSTRAINT FK_UserAsset_AssetType FOREIGN KEY (AssetTypeID) REFERENCES AssetType(AssetTypeID)
);

CREATE TABLE UserFinancialForecast (
    UserFinancialForecastID SERIAL PRIMARY KEY,
    UserFinancialForecastName VARCHAR(100) NOT NULL, 
    UserCategoryID INT, -- FK CustomerCategory
	EntityTypeName VARCHAR(50) NOT NULL, --Allowed Values: "User", "User Income" or "User Asset "
	EntityTypeID INT NOT NULL,			 -- If EntityTypeName "User" Then UserProfileID
										 -- If EntityTypeName "User Asset" Then UserAssetID 
									     -- If EntityTypeName "User Income" Then UserFinancialActualID
	EntityItemTypeName VARCHAR(50) NOT NULL, -- Allowed Values:
											 -- If EntityTypeName="User" then "Income" or "Expense"
											-- If EntityTypeName="User Income" then "Tax" or "Expense"
									        -- If EntityTypeName="User Asset" then "Income", "Expense" or "Tax"
    EntityItemTypeID INT NOT NULL,--If EntityItemTypeName="Income" Then IncomeTypeID
							      --If EntityItemTypeName="Expense" Then ExpenseTypeID
							      --If EntityItemTypeName="Tax" then TaxTypeID
    UserFinancialForecastAmount DECIMAL(15,2) NOT NULL,
    UserFinancialForecastBeginDate DATE NOT NULL,
    UserFinancialForecastEndDate DATE,
	CurrencyID INT NOT NULL, -- FK Currency
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserFinancialForecast_UserCategory FOREIGN KEY (UserCategoryID) REFERENCES UserCategory(UserCategoryID), 
	CONSTRAINT FK_UserFinancialForecast_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID)
);



-- UserFinancialActual
CREATE TABLE UserFinancialActual (
    UserFinancialActualID SERIAL PRIMARY KEY,
    UserFinancialActualtName VARCHAR(100) NOT NULL,
	UserCategoryID INT, -- FK CustomerCategory
	UserFinancialForecastID INT, -- FK UserFinancialForecast which Actual is connected if applicable (that is used to do the insight)
	EntityTypeName VARCHAR(50) NOT NULL, --Allowed Values: "User", "User Income" or "User Asset "
	EntityTypeID INT NOT NULL, -- If EntityTypeName "User" Then UserProfileID
							   -- If EntityTypeName"User Asset" Then UserAssetID 
							   -- If EntityTypeName"User Income" Then UserFinancialActualID
	EntityItemTypeName VARCHAR(50) NOT NULL, -- Allowed Values:
											 -- If Entity Name="User" then "Income or Expense"
									         -- If Entity Name="User Income" then "Tax or Expense"
									       -- If Entity Name="User Asset" then "Income, Expense or Tax"
    EntityItemTypeID INT NOT NULL,--If EntityItemTypeName="Income" Then IncomeTypeID
							      --If "Expense" Then ExpenseTypeID
							      --If "Tax" then TaxTypeID
    UserFinancialActualAmount DECIMAL(15,2) NOT NULL,
    UserFinancialActualBeginDate DATE NOT NULL,
    UserFinancialActualEndDate DATE,
	CurrencyID INT NOT NULL, -- FK Currency
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserFinancialActual_UserCategory FOREIGN KEY (UserCategoryID) REFERENCES UserCategory(UserCategoryID), 
	CONSTRAINT FK_UserFinancialActual_UserFinancialForecast FOREIGN KEY (UserFinancialForecastID) REFERENCES UserFinancialForecast(UserFinancialForecastID),
	CONSTRAINT FK_UserFinancialActual_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID)
);



-- Financial Insight Type Table
CREATE TABLE FinancialInsightType (
    FinancialInsightTypeID SERIAL PRIMARY KEY,
    InsightTypeName VARCHAR(100) NOT NULL UNIQUE, -- Nome do tipo de insight
    InsightDescription VARCHAR(255) NULL -- Descrição opcional
);

/*
Tipos de insights sugeridos:
1. Over Budget       - Quando o valor real excede o previsto
2. Under Budget      - Quando o valor real é menor que o previsto
3. On Target         - Quando o valor real está dentro de uma margem aceitável
4. High Variance     - Quando a diferença entre previsto e real é muito grande
5. Low Variance      - Quando a diferença entre previsto e real é pequena
6. Unexpected Expense - Quando um gasto não planejado é identificado
7. Missing Income    - Quando uma receita esperada não ocorreu
*/

-- Financial Insight table
CREATE TABLE FinancialInsight (
    FinancialInsightID SERIAL PRIMARY KEY,
	FinancialInsightTypeID INT NOT NULL, -- FK FinancialInsightTypeID
    FinancialForecastID INT NOT NULL, -- FK  FinancialForecast
    FinancialActualID INT NOT NULL, -- FK FinancialActual
    UserProfileID INT NOT NULL, -- FK  UserProfile
    UserCategoryID INT NOT NULL, -- FK UserCategory
    ItemTypeName VARCHAR(50) NOT NULL, -- Tipo (Ex.: Tax, Expense, Income)
    ForecastAmount DECIMAL(15,2) NOT NULL, -- Forecast Amount
    ActualAmount DECIMAL(15,2) NOT NULL, -- Actual Amount
    Variance DECIMAL(15,2) NOT NULL, -- Amount delta
    VariancePercentage DECIMAL(5,2) NOT NULL, -- delta %
    CalculatedDate DATE NOT NULL, -- calculation date
	CurrencYID INT NOT NULL,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_FinancialInsight_UserFinancialForecast FOREIGN KEY (FinancialForecastID) REFERENCES UserFinancialForecast(UserFinancialForecastID) ON DELETE CASCADE,
    CONSTRAINT FK_FinancialInsight_UserFinancialActual FOREIGN KEY (FinancialActualID) REFERENCES UserFinancialActual(UserFinancialActualID) ON DELETE CASCADE,
    CONSTRAINT FK_FinancialInsight_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID),
    CONSTRAINT FK_FinancialInsight_UserCategory FOREIGN KEY (UserCategoryID) REFERENCES UserCategory(UserCategoryID),
    CONSTRAINT FK_FinancialInsight_FinancialInsightType FOREIGN KEY (FinancialInsightTypeID) REFERENCES FinancialInsightType(FinancialInsightTypeID),
	CONSTRAINT FK_FinancialInsight_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID)

);