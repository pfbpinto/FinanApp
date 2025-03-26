
--------------------------------------------------------------------------------------------------
--------------------------------------REFERENCE TABLES---------------------------------------------
--------------------------------------------------------------------------------------------------



-- The entities which will be present in all references and user and groups line itens

CREATE TABLE Entity (
    EntityID SERIAL PRIMARY KEY,
    EntityName VARCHAR(150) NOT NULL,
    EntityType VARCHAR(150) NOT NULL,
	EntityCategory VARCHAR(10) NOT NULL,
	IsActive BOOLEAN NOT NULL DEFAULT TRUE,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

-- Income Type
CREATE TABLE IncomeType (
    IncomeTypeID SERIAL PRIMARY KEY,
    IncomeTypeName VARCHAR(100) NOT NULL,
    IncomeDescription TEXT NOT NULL,
	EntityID INT NOT NULL, --FK
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_IncomeType_Entity FOREIGN KEY (EntityID) REFERENCES Entity(EntityID) ON DELETE CASCADE
);


-- Tax Type
CREATE TABLE TaxType (
    TaxTypeID SERIAL PRIMARY KEY,
    TaxTypeName VARCHAR(100) NOT NULL,
    TaxDescription TEXT NOT NULL,
    TaxCountry VARCHAR(100) NOT NULL,
    TaxJurisdiction VARCHAR(100) NOT NULL,
    TaxPercentage DECIMAL(5,2),
	EntityID INT NOT NULL, --FK
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_TaxType_Entity FOREIGN KEY (EntityID) REFERENCES Entity(EntityID) ON DELETE CASCADE
);

-- Expense Type
CREATE TABLE ExpenseType (
    ExpenseTypeID SERIAL PRIMARY KEY,
    ExpenseTypeName VARCHAR(100) NOT NULL,
    ExpenseDescription TEXT NOT NULL,
    EntityID INT NOT NULL, --FK
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_ExpenseType_Entity FOREIGN KEY (EntityID) REFERENCES Entity(EntityID) ON DELETE CASCADE
);

-- Asset Type
CREATE TABLE AssetType (
    AssetTypeID SERIAL PRIMARY KEY,
    AssetTypeName VARCHAR(100) NOT NULL,
    AssetTypeDescription TEXT NOT NULL,
	EntityID INT NOT NULL, --FK
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_AssetType_Entity FOREIGN KEY (EntityID) REFERENCES Entity(EntityID) ON DELETE CASCADE
);

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
CREATE TABLE Recurrency (
    RecurrencyID SERIAL PRIMARY KEY,
	RecurrencyName VARCHAR(50) NOT NULL,
	RecurrencyPeriod VARCHAR(50) NOT NULL,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


--------------------------------------------------------------------------------------------------
--------------------------------USER MANAGEMENT TABLES--------------------------------------------
--------------------------------------------------------------------------------------------------
/* Here is all user related tables, even tough users are related to groups, they are completely idependent */


-- User Profile: Named as user profile as "User" is a native command in sql, to avoid issues, changed to userProfile
CREATE TABLE UserProfile (
    UserProfileID SERIAL PRIMARY KEY,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(255) NOT NULL,
	DateOfBirth DATE NOT NULL,
    UserPassword VARCHAR(150) NOT NULL,
    EmailAddress VARCHAR(255) UNIQUE NOT NULL,
    UserSubscription BIT NOT NULL, -- Boolean Active/inactive
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Asset Management. as Asset is a complex entity on it's own, will be manage separately. On MVP, assets will either be represented on one or other. If it's owned buy the user

CREATE TABLE UserAsset (
    UserAssetID SERIAL PRIMARY KEY,
    AssetTypeID INT NOT NULL, -- FK para a tabela AssetType
    UserProfileID INT NOT NULL, -- FK para a tabela UserProfile
    UserAssetName VARCHAR(100) NOT NULL, 
    UserAssetValueAmount DECIMAL(15,2) NOT NULL CHECK (UserAssetValueAmount >= 0),
    UserAssetAcquisitionBeginDate DATE NOT NULL,
    UserAssetAcquisitionEndDate DATE CHECK (UserAssetAcquisitionEndDate >= UserAssetAcquisitionBeginDate),
    IsActive BOOLEAN NOT NULL DEFAULT TRUE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserAsset_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID) ON DELETE CASCADE,
    CONSTRAINT FK_UserAsset_AssetType FOREIGN KEY (AssetTypeID) REFERENCES AssetType(AssetTypeID)
);



-- User Category
CREATE TABLE UserCategory (
    UserCategoryID SERIAL PRIMARY KEY,
    UserCategoryName VARCHAR(255) NOT NULL,
    UserProfileID INT NOT NULL,
	EntityID INT NOT NULL , -- FK
    FinancialGroupEntityItemID INT,--If EntityItemTypeName="Tax" then TaxTypeID from TaxType Table
							 --If EntityItemTypeName="Income" then IncomeTypeID from IncomeType Table
							 --If EntityItemTypeName="Expense" then ExpenseTypeID from ExpenseType Table
	IsActive BOOLEAN NOT NULL DEFAULT TRUE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserCategory_UserProfile FOREIGN KEY (UserProfileID) REFERENCES UserProfile(UserProfileID) ON DELETE CASCADE,
	CONSTRAINT FK_UserCategory_Entity FOREIGN KEY (EntityID) REFERENCES Entity(EntityID) ON DELETE CASCADE
);


-- Financial User Item -- This is the custom name users can give to a Income, Expense, Asset, Asset Income, or Asset Expense
CREATE TABLE FinancialUserItem (
    FinancialUserItemID SERIAL PRIMARY KEY,
    FinancialUserItemName VARCHAR(255) NOT NULL,
    EntityID INT NOT NULL , -- FK - From the entity you also define if it is User or Asset and what type of Record this is (Parent or Child), if child, it must have the ParentFinancialUserItem populated
	UserEntityID INT NOT NULL, -- Here is either the UserProfileID or AssetID, depending on the entity selected
	RecurrencyID INT NOT NULL, --FK
    FinancialUserEntityItemID INT,--If EntityItemTypeName="Tax" then TaxTypeID from TaxType Table
							 --If EntityItemTypeName="Income" then IncomeTypeID from IncomeType Table
							 --If EntityItemTypeName="Expense" then ExpenseTypeID from ExpenseType Table
	ParentFinancialUserItemID INT, -- This will reference cases where the item is a child. (e.g.: A "User Income" has a "Tax" as child, the Tax User Item would refer here to which "Income" it's tied to)
	IsActive BOOLEAN NOT NULL DEFAULT TRUE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_FinancialUserItem_Entity FOREIGN KEY (EntityID) REFERENCES Entity(EntityID) ON DELETE CASCADE,
	CONSTRAINT FK_FinancialUserItem_Recurrency FOREIGN KEY (RecurrencyID) REFERENCES Recurrency(RecurrencyID) ON DELETE CASCADE,
	CONSTRAINT FK_FinancialUserItem_ParentFinancialUserItemID FOREIGN KEY (ParentFinancialUserItemID) REFERENCES FinancialUserItem(FinancialUserItemID)
);


	
 -- User Forecast Table
 
CREATE TABLE UserFinancialForecast (
    UserFinancialForecastID SERIAL PRIMARY KEY,
    UserCategoryID INT, -- FK CustomerCategory
	FinancialUserItemID INT NOT NULL, -- FK Financial User Item
    UserFinancialForecastBeginDate DATE NOT NULL,
    UserFinancialForecastEndDate DATE,
	UserFinancialForecastAmount DECIMAL(15,2) NOT NULL, -- Forecast Amount
	CurrencyID INT NOT NULL, -- FK Currency
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserFinancialForecast_UserCategory FOREIGN KEY (UserCategoryID) REFERENCES UserCategory(UserCategoryID), 
	CONSTRAINT FK_UserFinancialForecast_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID),
	CONSTRAINT FK_UserFinancialForecast_FinancialUserItemID FOREIGN KEY (FinancialUserItemID) REFERENCES FinancialUserItem(FinancialUserItemID)

);

-- UserFinancialActual
CREATE TABLE UserFinancialActual (
    UserFinancialActualID SERIAL PRIMARY KEY,
    UserCategoryID INT, -- FK CustomerCategory
	FinancialUserItemID INT NOT NULL, -- FK Financial User Item
    UserFinancialActualtBeginDate DATE NOT NULL,
    UserFinancialActualEndDate DATE,
	UserFinancialActualAmount DECIMAL(15,2) NOT NULL, -- Forecast Amount
	CurrencyID INT NOT NULL, -- FK Currency
	Note TEXT,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_UserFinancialActual_UserCategory FOREIGN KEY (UserCategoryID) REFERENCES UserCategory(UserCategoryID), 
	CONSTRAINT FK_UserFinancialActual_Currency FOREIGN KEY (CurrencyID) REFERENCES Currency(CurrencyID)
);

-- Financial User Forecast and Actuals relationship table
CREATE TABLE UserForecastActualRelation(
    UserForecastActualRelationID SERIAL PRIMARY KEY,
	UserFinancialActualID INT NOT NULL, -- FK
	UserFinancialForecastID INT NOT NULL, -- FK
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_UserForecastActualRelation_UserFinancialActual FOREIGN KEY (UserFinancialActualID) REFERENCES UserFinancialActual(UserFinancialActualID),
	CONSTRAINT FK_UserForecastActualRelation_UserFinancialForecast FOREIGN KEY (UserFinancialForecastID) REFERENCES UserFinancialForecast(UserFinancialForecastID)

);